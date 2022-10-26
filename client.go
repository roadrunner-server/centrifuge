package centrifuge

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	v1Client "go.buf.build/grpc/go/roadrunner-server/api/proto/centrifugo/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	// Will register via init
	_ "google.golang.org/grpc/encoding/gzip"
)

type client struct {
	mu sync.RWMutex

	log      *zap.Logger
	addr     string
	tls      *TLS
	compress bool

	centrifugoClient v1Client.CentrifugoApiClient
}

func newClient(addr string, tls *TLS, log *zap.Logger, compress bool) *client {
	return &client{
		addr:     addr,
		tls:      tls,
		compress: compress,
		log:      log,
	}
}

func (c *client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cb := backoff.NewConstantBackOff(time.Minute)

	opts := make([]grpc.CallOption, 0, 1)

	if c.compress {
		opts = append(opts, grpc.UseCompressor("gzip"))
	}

	operation := func() error {
		if c.tls != nil {
			cert, err := tls.LoadX509KeyPair(c.tls.Cert, c.tls.Key)
			if err != nil {
				return err
			}

			tlscfg := &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			}

			conn, err := grpc.Dial(c.addr, grpc.WithDefaultCallOptions(opts...), grpc.WithTransportCredentials(credentials.NewTLS(tlscfg)))
			if err != nil {
				c.log.Debug("attempted to connect to the centrifugo server with TLS, retrying", zap.Error(err))
				return err
			}

			c.centrifugoClient = v1Client.NewCentrifugoApiClient(conn)
			return nil
		}

		// non-tls
		conn, err := grpc.Dial(c.addr, grpc.WithDefaultCallOptions(opts...), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			c.log.Debug("attempted to connect to the centrifugo server, retrying", zap.Error(err))
			return err
		}

		c.centrifugoClient = v1Client.NewCentrifugoApiClient(conn)
		return nil
	}

	err := backoff.Retry(operation, cb)
	if err != nil {
		return err
	}

	c.log.Debug("connected to the centrifugo server")

	return nil
}

func (c *client) client() v1Client.CentrifugoApiClient {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.centrifugoClient
}
