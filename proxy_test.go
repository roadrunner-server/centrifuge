package centrifuge

import (
	"context"
	"errors"
	"sync"
	"testing"

	centrifugov1 "github.com/roadrunner-server/api-go/v6/centrifugo/proxy/v1"
	"github.com/roadrunner-server/pool/v2/payload"
	staticPool "github.com/roadrunner-server/pool/v2/pool/static_pool"
	"github.com/roadrunner-server/pool/v2/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// fakePool implements the Pool interface. Exec always returns an error so each
// proxy method exits right after the gRPC-metadata guard block (the lines under
// test); staticPool.PExec cannot be constructed outside its package, so a
// success path is not reproducible here.
type fakePool struct {
	execErr error
}

func (f *fakePool) Workers() []*worker.Process           { return nil }
func (f *fakePool) RemoveWorker(_ context.Context) error { return nil }
func (f *fakePool) AddWorker() error                     { return nil }
func (f *fakePool) Reset(_ context.Context) error        { return nil }
func (f *fakePool) Destroy(_ context.Context)            {}

func (f *fakePool) Exec(_ context.Context, _ *payload.Payload, _ chan struct{}) (chan *staticPool.PExec, error) {
	return nil, f.execErr
}

func newTestProxy() *Proxy {
	return &Proxy{
		log: testLogger(),
		pw:  newPoolMuWrapper(&fakePool{execErr: errors.New("exec failed")}, &sync.RWMutex{}),
	}
}

// TestProxyMethodsMetadataGuard exercises every proxy handler with and without
// incoming gRPC metadata. The no-metadata case is the regression guard: the old
// discarded-ok pattern left md as a nil map, so md.Append panicked on every
// unauthenticated call.
func TestProxyMethodsMetadataGuard(t *testing.T) {
	p := newTestProxy()

	calls := map[string]func(context.Context) error{
		"Connect": func(ctx context.Context) error { _, err := p.Connect(ctx, &centrifugov1.ConnectRequest{}); return err },
		"Refresh": func(ctx context.Context) error { _, err := p.Refresh(ctx, &centrifugov1.RefreshRequest{}); return err },
		"Subscribe": func(ctx context.Context) error {
			_, err := p.Subscribe(ctx, &centrifugov1.SubscribeRequest{})
			return err
		},
		"Publish": func(ctx context.Context) error { _, err := p.Publish(ctx, &centrifugov1.PublishRequest{}); return err },
		"RPC":     func(ctx context.Context) error { _, err := p.RPC(ctx, &centrifugov1.RPCRequest{}); return err },
		"SubRefresh": func(ctx context.Context) error {
			_, err := p.SubRefresh(ctx, &centrifugov1.SubRefreshRequest{})
			return err
		},
		"NotifyCacheEmpty": func(ctx context.Context) error {
			_, err := p.NotifyCacheEmpty(ctx, &centrifugov1.NotifyCacheEmptyRequest{})
			return err
		},
		"NotifyChannelState": func(ctx context.Context) error {
			_, err := p.NotifyChannelState(ctx, &centrifugov1.NotifyChannelStateRequest{})
			return err
		},
	}

	ctxNoMD := t.Context()
	ctxWithMD := metadata.NewIncomingContext(t.Context(), metadata.Pairs("authorization", "bearer x"))

	for name, call := range calls {
		var err error

		require.NotPanics(t, func() { err = call(ctxNoMD) }, "%s panicked without incoming metadata", name)
		require.Error(t, err, "%s: expected error from failing Exec (no metadata)", name)

		require.NotPanics(t, func() { err = call(ctxWithMD) }, "%s panicked with incoming metadata", name)
		require.Error(t, err, "%s: expected error from failing Exec (with metadata)", name)
	}
}

func TestProxyStreamingNotSupported(t *testing.T) {
	p := newTestProxy()

	require.Error(t, p.SubscribeUnidirectional(&centrifugov1.SubscribeRequest{}, nil))
	require.Error(t, p.SubscribeBidirectional(nil))
}
