package centrifuge

import (
	"context"

	"github.com/goccy/go-json"
	centrifugov1 "github.com/roadrunner-server/api/v4/build/centrifugo/proxy/v1"
	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/goridge/v3/pkg/frame"
	"github.com/roadrunner-server/pool/payload"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Proxy struct {
	log *zap.Logger
	pw  *wrapper
}

func (p *Proxy) Connect(ctx context.Context, request *centrifugov1.ConnectRequest) (*centrifugov1.ConnectResponse, error) {
	p.log.Debug("got connect proxy request")
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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	cr := &centrifugov1.ConnectResponse{}
	err = proto.Unmarshal(re.Body, cr)

	if err != nil {
		return nil, err
	}

	p.log.Debug("finished connect proxy request")
	return cr, nil
}

func (p *Proxy) Refresh(ctx context.Context, request *centrifugov1.RefreshRequest) (*centrifugov1.RefreshResponse, error) {
	p.log.Debug("got refresh proxy request")
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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	rr := &centrifugov1.RefreshResponse{}

	err = proto.Unmarshal(re.Body, rr)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished refresh proxy request")
	return rr, nil
}

func (p *Proxy) Subscribe(ctx context.Context, request *centrifugov1.SubscribeRequest) (*centrifugov1.SubscribeResponse, error) {
	p.log.Debug("got subscribe proxy request")
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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	sr := &centrifugov1.SubscribeResponse{}

	err = proto.Unmarshal(re.Body, sr)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished subscribe proxy request")
	return sr, nil
}

func (p *Proxy) Publish(ctx context.Context, request *centrifugov1.PublishRequest) (*centrifugov1.PublishResponse, error) {
	p.log.Debug("got publish proxy request")
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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	pr := &centrifugov1.PublishResponse{}

	err = proto.Unmarshal(re.Body, pr)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished publish proxy request")
	return pr, nil
}

func (p *Proxy) RPC(ctx context.Context, request *centrifugov1.RPCRequest) (*centrifugov1.RPCResponse, error) {
	p.log.Debug("got RPC proxy request", zap.String("method", request.Method))
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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	rresp := &centrifugov1.RPCResponse{}

	err = proto.Unmarshal(re.Body, rresp)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished RPC proxy request")
	return rresp, nil
}

func (p *Proxy) SubRefresh(ctx context.Context, request *centrifugov1.SubRefreshRequest) (*centrifugov1.SubRefreshResponse, error) {
	p.log.Debug("got RPC SubRefresh request", zap.String("channel", request.Channel))

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

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	rresp := &centrifugov1.SubRefreshResponse{}

	err = proto.Unmarshal(re.Body, rresp)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished RPC SubRefresh request")
	return rresp, nil
}

func (p *Proxy) NotifyChannelState(ctx context.Context, request *centrifugov1.NotifyChannelStateRequest) (*centrifugov1.NotifyChannelStateResponse, error) {
	p.log.Debug("got NotifyChannelState request")

	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromIncomingContext(ctx)
	md.Append("type", "notifychannelstate")

	meta, err := json.Marshal(md)
	if err != nil {
		return nil, err
	}

	pld := &payload.Payload{
		Context: meta,
		Body:    data,
		Codec:   frame.CodecProto,
	}

	re, err := p.pw.Exec(ctx, pld)
	if err != nil {
		return nil, err
	}

	rresp := &centrifugov1.NotifyChannelStateResponse{}

	err = proto.Unmarshal(re.Body, rresp)
	if err != nil {
		return nil, err
	}

	p.log.Debug("finished NotifyChannelState request")
	return rresp, nil
}

func (p *Proxy) SubscribeUnidirectional(_ *centrifugov1.SubscribeRequest, _ centrifugov1.CentrifugoProxy_SubscribeUnidirectionalServer) error {
	p.log.Debug("got SubscribeUnidirectional request")

	return errors.Str("not supported")
}

func (p *Proxy) SubscribeBidirectional(_ centrifugov1.CentrifugoProxy_SubscribeBidirectionalServer) error {
	p.log.Debug("got StreamSubRequest request")

	return errors.Str("not supported")
}
