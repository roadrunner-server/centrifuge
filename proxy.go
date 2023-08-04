package centrifuge

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/roadrunner-server/goridge/v3/pkg/frame"
	"github.com/roadrunner-server/sdk/v4/payload"
	centrifugov1 "github.com/roadrunner-server/api/v4/build/centrifugo/proxy/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Proxy struct {
	p *Plugin
}

func (p *Proxy) Connect(ctx context.Context, request *centrifugov1.ConnectRequest) (*centrifugov1.ConnectResponse, error) {
	p.p.log.Debug("got connect proxy request")
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "connect")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}
	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	cr := &centrifugov1.ConnectResponse{}

	err = proto.Unmarshal(resp.Body(), cr)
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func (p *Proxy) Refresh(ctx context.Context, request *centrifugov1.RefreshRequest) (*centrifugov1.RefreshResponse, error) {
	p.p.log.Debug("got refresh proxy request")
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "refresh")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	if resp.Payload().IsStream {
		sc <- struct{}{}
		return nil, errors.New("streaming response not supported")
	}

	rr := &centrifugov1.RefreshResponse{}

	err = proto.Unmarshal(resp.Body(), rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (p *Proxy) Subscribe(ctx context.Context, request *centrifugov1.SubscribeRequest) (*centrifugov1.SubscribeResponse, error) {
	p.p.log.Debug("got subscribe proxy request")
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "subscribe")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	if resp.Payload().IsStream {
		sc <- struct{}{}
		return nil, errors.New("streaming response not supported")
	}

	sr := &centrifugov1.SubscribeResponse{}

	err = proto.Unmarshal(resp.Body(), sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func (p *Proxy) Publish(ctx context.Context, request *centrifugov1.PublishRequest) (*centrifugov1.PublishResponse, error) {
	p.p.log.Debug("got publish proxy request")
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "publish")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	if resp.Payload().IsStream {
		sc <- struct{}{}
		return nil, errors.New("streaming response not supported")
	}

	pr := &centrifugov1.PublishResponse{}

	err = proto.Unmarshal(resp.Body(), pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *Proxy) RPC(ctx context.Context, request *centrifugov1.RPCRequest) (*centrifugov1.RPCResponse, error) {
	p.p.log.Debug("got RPC proxy request", zap.String("method", request.Method))
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "rpc")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	if resp.Payload().IsStream {
		sc <- struct{}{}
		return nil, errors.New("streaming response not supported")
	}

	rresp := &centrifugov1.RPCResponse{}

	err = proto.Unmarshal(resp.Body(), rresp)
	if err != nil {
		return nil, err
	}

	return rresp, nil
}

func (p *Proxy) SubRefresh(ctx context.Context, request *centrifugov1.SubRefreshRequest) (*centrifugov1.SubRefreshResponse, error) {
	p.p.log.Debug("got RPC SubRefresh request", zap.String("channel", request.Channel))

	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "subrefresh")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	p.p.mu.RLock()
	sc := make(chan struct{}, 1)
	re, err := p.p.pool.Exec(ctx, pld, sc)
	p.p.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	resp := <-re
	if resp.Payload().IsStream {
		sc <- struct{}{}
		return nil, errors.New("streaming response not supported")
	}
	rresp := &centrifugov1.SubRefreshResponse{}

	err = proto.Unmarshal(resp.Body(), rresp)
	if err != nil {
		return nil, err
	}

	return rresp, nil
}
