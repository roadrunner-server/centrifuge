package centrifuge

import (
	"context"

	"github.com/roadrunner-server/goridge/v3/pkg/frame"
	"github.com/roadrunner-server/sdk/v3/payload"
	"github.com/segmentio/encoding/proto"
	centrifugov1 "go.buf.build/grpc/go/roadrunner-server/api/proto/centrifugo/proxy/v1"
)

type Proxy struct {
	p *Plugin
}

func (p *Proxy) Connect(ctx context.Context, request *centrifugov1.ConnectRequest) (*centrifugov1.ConnectResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: []byte("connect"),
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	resp, err := p.p.pool.Exec(ctx, pld)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	cr := &centrifugov1.ConnectResponse{}

	err = proto.Unmarshal(resp.Body, &cr)
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (p *Proxy) Refresh(ctx context.Context, request *centrifugov1.RefreshRequest) (*centrifugov1.RefreshResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: []byte("refresh"),
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	resp, err := p.p.pool.Exec(ctx, pld)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	rr := &centrifugov1.RefreshResponse{}

	err = proto.Unmarshal(resp.Body, &rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (p *Proxy) Subscribe(ctx context.Context, request *centrifugov1.SubscribeRequest) (*centrifugov1.SubscribeResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: []byte("subscribe"),
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	resp, err := p.p.pool.Exec(ctx, pld)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	sr := &centrifugov1.SubscribeResponse{}

	err = proto.Unmarshal(resp.Body, &sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func (p *Proxy) Publish(ctx context.Context, request *centrifugov1.PublishRequest) (*centrifugov1.PublishResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: []byte("publish"),
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	resp, err := p.p.pool.Exec(ctx, pld)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	pr := &centrifugov1.PublishResponse{}

	err = proto.Unmarshal(resp.Body, &pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *Proxy) RPC(ctx context.Context, request *centrifugov1.RPCRequest) (*centrifugov1.RPCResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: []byte("rpc"),
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	resp, err := p.p.pool.Exec(ctx, pld)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	rresp := &centrifugov1.RPCResponse{}

	err = proto.Unmarshal(resp.Body, &rresp)
	if err != nil {
		return nil, err
	}

	return rresp, nil
}
