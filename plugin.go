package centrifuge

import (
	"context"
	stderr "errors"
	"sync"

	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/sdk/v3/payload"
	"github.com/roadrunner-server/sdk/v3/pool"
	staticPool "github.com/roadrunner-server/sdk/v3/pool/static_pool"
	"github.com/roadrunner-server/sdk/v3/utils"
	"github.com/roadrunner-server/sdk/v3/worker"
	centrifugov1 "go.buf.build/grpc/go/roadrunner-server/api/proto/centrifugo/proxy/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	name string = "centrifuge"
)

type Configurer interface {
	// UnmarshalKey takes a single key and unmarshal it into a Struct.
	UnmarshalKey(name string, out any) error
	// Has checks if config section exists.
	Has(name string) bool
}

type Pool interface {
	// Workers returns worker list associated with the pool.
	Workers() (workers []*worker.Process)

	// Exec payload
	Exec(ctx context.Context, p *payload.Payload) (*payload.Payload, error)

	// Reset kill all workers inside the watcher and replaces with new
	Reset(ctx context.Context) error

	// Destroy all underlying stack (but let them complete the task).
	Destroy(ctx context.Context)
}

// Server creates workers for the application.
type Server interface {
	NewPool(ctx context.Context, cfg *pool.Config, env map[string]string, _ *zap.Logger) (*staticPool.Pool, error)
}

type Plugin struct {
	mu  sync.RWMutex
	cfg *Config

	log    *zap.Logger
	server Server

	pool Pool
}

func (p *Plugin) Init(cfg Configurer, log *zap.Logger, server Server) error {
	const op = errors.Op("centrifuge_plugin_init")
	if !cfg.Has(name) {
		return errors.E(op, errors.Disabled)
	}

	err := cfg.UnmarshalKey(name, &p.cfg)
	if err != nil {
		return errors.E(op, err)
	}

	p.cfg.InitDefaults()

	p.log = new(zap.Logger)
	*p.log = *log
	p.server = server

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)
	const op = errors.Op("centrifugo_serve")

	var err error
	p.mu.Lock()
	p.pool, err = p.server.NewPool(context.Background(), p.cfg.Pool, nil, nil)
	p.mu.Unlock()
	if err != nil {
		errCh <- err
		return errCh
	}

	l, err := utils.CreateListener(p.cfg.ProxyAddress)
	if err != nil {
		errCh <- errors.E(op, err)
		return errCh
	}

	server := grpc.NewServer()
	centrifugov1.RegisterCentrifugoProxyServer(server, &Proxy{p: p})

	go func() {
		errL := server.Serve(l)
		if errL != nil {
			if stderr.Is(errL, grpc.ErrServerStopped) {
				p.log.Info("grpc proxy stopped")
				return
			}

			p.log.Error("grpc proxy error", zap.Error(errL))
		}
	}()

	return errCh
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) Name() string {
	return name
}
