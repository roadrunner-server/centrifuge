package centrifuge

import (
	"context"
	stderr "errors"
	"sync"
	"time"

	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/sdk/v3/metrics"
	"github.com/roadrunner-server/sdk/v3/payload"
	"github.com/roadrunner-server/sdk/v3/pool"
	staticPool "github.com/roadrunner-server/sdk/v3/pool/static_pool"
	"github.com/roadrunner-server/sdk/v3/state/process"
	"github.com/roadrunner-server/sdk/v3/utils"
	"github.com/roadrunner-server/sdk/v3/worker"
	centrifugov1 "go.buf.build/grpc/go/roadrunner-server/api/centrifugo/proxy/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	name string = "centrifuge"

	RRMode           = "RR_MODE"
	RRModeCentrifuge = "centrifuge"
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

	// proxy server
	gRPCServer    *grpc.Server
	client        *client
	statsExporter *metrics.StatsExporter

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

	err = p.cfg.InitDefaults()
	if err != nil {
		return err
	}

	p.log = new(zap.Logger)
	*p.log = *log
	p.server = server
	p.gRPCServer = grpc.NewServer()
	p.client = newClient(p.cfg.GrpcApiAddress, p.cfg.TLS, p.log, p.cfg.UseCompressor)
	p.statsExporter = newWorkersExporter(p)

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)
	const op = errors.Op("centrifuge_serve")

	var err error
	p.mu.Lock()
	p.pool, err = p.server.NewPool(context.Background(), p.cfg.Pool, map[string]string{RRMode: RRModeCentrifuge}, nil)
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

	p.mu.Lock()
	centrifugov1.RegisterCentrifugoProxyServer(p.gRPCServer, &Proxy{p: p})
	p.mu.Unlock()

	go func() {
		errL := p.gRPCServer.Serve(l)
		if errL != nil {
			if stderr.Is(errL, grpc.ErrServerStopped) {
				p.log.Info("grpc proxy stopped")
				return
			}

			p.log.Error("grpc proxy error", zap.Error(errL))
		}
	}()

	err = p.client.connect()
	if err != nil {
		errCh <- err
		return errCh
	}

	return errCh
}

func (p *Plugin) Stop() error {
	p.mu.Lock()
	p.gRPCServer.Stop()
	p.mu.Unlock()
	return nil
}

// Workers returns slice with the process states for the workers
func (p *Plugin) Workers() []*process.State {
	p.mu.RLock()
	defer p.mu.RUnlock()

	workers := p.workers()
	if workers == nil {
		return nil
	}

	ps := make([]*process.State, 0, len(workers))
	for i := 0; i < len(workers); i++ {
		state, err := process.WorkerProcessState(workers[i])
		if err != nil {
			return nil
		}
		ps = append(ps, state)
	}

	return ps
}

// Reset destroys the old pool and replaces it with new one, waiting for old pool to die
func (p *Plugin) Reset() error {
	const op = errors.Op("centrifuge_plugin_reset")

	p.log.Info("reset signal was received")

	ctxTout, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	if p.pool == nil {
		p.log.Info("pool is nil, nothing to reset")
		return nil
	}

	err := p.pool.Reset(ctxTout)
	if err != nil {
		return errors.E(op, err)
	}

	p.log.Info("plugin was successfully reset")
	return nil
}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) RPC() any {
	return &rpc{
		client: p.client,
		log:    p.log,
	}
}

// internal
func (p *Plugin) workers() []*worker.Process {
	if p == nil || p.pool == nil {
		return nil
	}
	return p.pool.Workers()
}
