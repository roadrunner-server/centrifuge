package centrifuge

import (
	"context"
	"errors"
	"sync"

	"github.com/roadrunner-server/goridge/v3/pkg/frame"
	"github.com/roadrunner-server/pool/payload"
)

type wrapper struct {
	stopChPool sync.Pool
	mu         *sync.RWMutex
	pool       Pool
}

func newPoolMuWrapper(pool Pool, mu *sync.RWMutex) *wrapper {
	return &wrapper{
		stopChPool: sync.Pool{
			New: func() any {
				return make(chan struct{}, 1)
			},
		},
		mu:   mu,
		pool: pool,
	}
}

func (p *wrapper) Exec(ctx context.Context, pld *payload.Payload) (*payload.Payload, error) {
	p.mu.RLock()
	sc := p.getStopCh()
	re, err := p.pool.Exec(ctx, pld, sc)
	p.mu.RUnlock()
	if err != nil {
		p.putStopCh(sc)

		return nil, err
	}

	var resp *payload.Payload
	select {
	case pl := <-re:
		if pl.Error() != nil {
			p.putStopCh(sc)

			return nil, pl.Error()
		}
		// streaming is not supported
		if pl.Payload().Flags&frame.STREAM != 0 {
			// stop the stream, do not return the channel
			sc <- struct{}{}

			return nil, errors.New("streaming response not supported")
		}

		// assign the payload
		resp = pl.Payload()
	default:
		p.putStopCh(sc)

		return nil, errors.New("worker empty response")
	}

	p.putStopCh(sc)

	return resp, nil
}

func (p *wrapper) getStopCh() chan struct{} {
	return p.stopChPool.Get().(chan struct{})
}

func (p *wrapper) putStopCh(ch chan struct{}) {
	select {
	case <-ch:
	default:
	}
	p.stopChPool.Put(ch)
}
