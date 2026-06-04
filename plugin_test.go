package centrifuge

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// testLogger returns a no-op slog logger shared across the package's unit tests.
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestPluginStopNilPool(t *testing.T) {
	// pool is nil when Serve returned early; Stop must not panic on Destroy.
	p := &Plugin{gRPCServer: grpc.NewServer()}

	var err error
	require.NotPanics(t, func() { err = p.Stop(t.Context()) })
	require.NoError(t, err)
}

func TestPluginWorkersNilPool(t *testing.T) {
	p := &Plugin{}

	assert.Nil(t, p.Workers())
}

func TestPluginResetNilPool(t *testing.T) {
	p := &Plugin{log: testLogger()}

	require.NoError(t, p.Reset())
}
