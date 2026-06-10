package centrifuge

import (
	"context"
	stderr "errors"
	"log/slog"

	"connectrpc.com/connect"
	v1Client "github.com/roadrunner-server/api-go/v6/centrifugo/api/v1"
	"github.com/roadrunner-server/api-go/v6/centrifugo/api/v1/apiV1connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Compile-time check that rpc implements the generated handler interface.
var _ apiV1connect.CentrifugoApiHandler = (*rpc)(nil)

// rpc exposes the Centrifugo server API (centrifugal.centrifugo.api.CentrifugoApi)
// on the RoadRunner Connect-RPC plane; every call is proxied 1:1 to the
// connected Centrifugo server.
type rpc struct {
	client *client
	log    *slog.Logger
}

// forward proxies a single CentrifugoApi call to the connected Centrifugo
// client, translating gRPC status errors into connect errors (the two code
// registries are aligned). The client is re-fetched on every call because
// connect() replaces it after reconnects.
func forward[Req, Resp any](
	ctx context.Context,
	r *rpc,
	req *connect.Request[Req],
	call func(v1Client.CentrifugoApiClient, context.Context, *Req, ...grpc.CallOption) (*Resp, error),
) (*connect.Response[Resp], error) {
	r.log.Debug("got api request", "procedure", req.Spec().Procedure)

	client := r.client.client()
	if client == nil {
		return nil, connect.NewError(connect.CodeUnavailable, stderr.New("RoadRunner is not ready yet, try in a few seconds"))
	}

	resp, err := call(client, ctx, req.Msg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, connect.NewError(connect.Code(st.Code()), stderr.New(st.Message()))
		}
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (r *rpc) Batch(ctx context.Context, req *connect.Request[v1Client.BatchRequest]) (*connect.Response[v1Client.BatchResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Batch)
}

func (r *rpc) Publish(ctx context.Context, req *connect.Request[v1Client.PublishRequest]) (*connect.Response[v1Client.PublishResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Publish)
}

func (r *rpc) Broadcast(ctx context.Context, req *connect.Request[v1Client.BroadcastRequest]) (*connect.Response[v1Client.BroadcastResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Broadcast)
}

func (r *rpc) Subscribe(ctx context.Context, req *connect.Request[v1Client.SubscribeRequest]) (*connect.Response[v1Client.SubscribeResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Subscribe)
}

func (r *rpc) Unsubscribe(ctx context.Context, req *connect.Request[v1Client.UnsubscribeRequest]) (*connect.Response[v1Client.UnsubscribeResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Unsubscribe)
}

func (r *rpc) Disconnect(ctx context.Context, req *connect.Request[v1Client.DisconnectRequest]) (*connect.Response[v1Client.DisconnectResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Disconnect)
}

func (r *rpc) Presence(ctx context.Context, req *connect.Request[v1Client.PresenceRequest]) (*connect.Response[v1Client.PresenceResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Presence)
}

func (r *rpc) PresenceStats(ctx context.Context, req *connect.Request[v1Client.PresenceStatsRequest]) (*connect.Response[v1Client.PresenceStatsResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.PresenceStats)
}

func (r *rpc) History(ctx context.Context, req *connect.Request[v1Client.HistoryRequest]) (*connect.Response[v1Client.HistoryResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.History)
}

func (r *rpc) HistoryRemove(ctx context.Context, req *connect.Request[v1Client.HistoryRemoveRequest]) (*connect.Response[v1Client.HistoryRemoveResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.HistoryRemove)
}

func (r *rpc) Info(ctx context.Context, req *connect.Request[v1Client.InfoRequest]) (*connect.Response[v1Client.InfoResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Info)
}

func (r *rpc) RPC(ctx context.Context, req *connect.Request[v1Client.RPCRequest]) (*connect.Response[v1Client.RPCResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.RPC)
}

func (r *rpc) Refresh(ctx context.Context, req *connect.Request[v1Client.RefreshRequest]) (*connect.Response[v1Client.RefreshResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Refresh)
}

func (r *rpc) Channels(ctx context.Context, req *connect.Request[v1Client.ChannelsRequest]) (*connect.Response[v1Client.ChannelsResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Channels)
}

func (r *rpc) Connections(ctx context.Context, req *connect.Request[v1Client.ConnectionsRequest]) (*connect.Response[v1Client.ConnectionsResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.Connections)
}

func (r *rpc) UpdateUserStatus(ctx context.Context, req *connect.Request[v1Client.UpdateUserStatusRequest]) (*connect.Response[v1Client.UpdateUserStatusResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.UpdateUserStatus)
}

func (r *rpc) GetUserStatus(ctx context.Context, req *connect.Request[v1Client.GetUserStatusRequest]) (*connect.Response[v1Client.GetUserStatusResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.GetUserStatus)
}

func (r *rpc) DeleteUserStatus(ctx context.Context, req *connect.Request[v1Client.DeleteUserStatusRequest]) (*connect.Response[v1Client.DeleteUserStatusResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeleteUserStatus)
}

func (r *rpc) BlockUser(ctx context.Context, req *connect.Request[v1Client.BlockUserRequest]) (*connect.Response[v1Client.BlockUserResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.BlockUser)
}

func (r *rpc) UnblockUser(ctx context.Context, req *connect.Request[v1Client.UnblockUserRequest]) (*connect.Response[v1Client.UnblockUserResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.UnblockUser)
}

func (r *rpc) RevokeToken(ctx context.Context, req *connect.Request[v1Client.RevokeTokenRequest]) (*connect.Response[v1Client.RevokeTokenResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.RevokeToken)
}

func (r *rpc) InvalidateUserTokens(ctx context.Context, req *connect.Request[v1Client.InvalidateUserTokensRequest]) (*connect.Response[v1Client.InvalidateUserTokensResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.InvalidateUserTokens)
}

func (r *rpc) DeviceRegister(ctx context.Context, req *connect.Request[v1Client.DeviceRegisterRequest]) (*connect.Response[v1Client.DeviceRegisterResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceRegister)
}

func (r *rpc) DeviceUpdate(ctx context.Context, req *connect.Request[v1Client.DeviceUpdateRequest]) (*connect.Response[v1Client.DeviceUpdateResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceUpdate)
}

func (r *rpc) DeviceRemove(ctx context.Context, req *connect.Request[v1Client.DeviceRemoveRequest]) (*connect.Response[v1Client.DeviceRemoveResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceRemove)
}

func (r *rpc) DeviceList(ctx context.Context, req *connect.Request[v1Client.DeviceListRequest]) (*connect.Response[v1Client.DeviceListResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceList)
}

func (r *rpc) DeviceTopicList(ctx context.Context, req *connect.Request[v1Client.DeviceTopicListRequest]) (*connect.Response[v1Client.DeviceTopicListResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceTopicList)
}

func (r *rpc) DeviceTopicUpdate(ctx context.Context, req *connect.Request[v1Client.DeviceTopicUpdateRequest]) (*connect.Response[v1Client.DeviceTopicUpdateResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.DeviceTopicUpdate)
}

func (r *rpc) UserTopicList(ctx context.Context, req *connect.Request[v1Client.UserTopicListRequest]) (*connect.Response[v1Client.UserTopicListResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.UserTopicList)
}

func (r *rpc) UserTopicUpdate(ctx context.Context, req *connect.Request[v1Client.UserTopicUpdateRequest]) (*connect.Response[v1Client.UserTopicUpdateResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.UserTopicUpdate)
}

func (r *rpc) SendPushNotification(ctx context.Context, req *connect.Request[v1Client.SendPushNotificationRequest]) (*connect.Response[v1Client.SendPushNotificationResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.SendPushNotification)
}

func (r *rpc) UpdatePushStatus(ctx context.Context, req *connect.Request[v1Client.UpdatePushStatusRequest]) (*connect.Response[v1Client.UpdatePushStatusResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.UpdatePushStatus)
}

func (r *rpc) CancelPush(ctx context.Context, req *connect.Request[v1Client.CancelPushRequest]) (*connect.Response[v1Client.CancelPushResponse], error) {
	return forward(ctx, r, req, v1Client.CentrifugoApiClient.CancelPush)
}
