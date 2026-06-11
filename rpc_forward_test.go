package centrifuge

import (
	"context"
	"net"
	"testing"

	"connectrpc.com/connect"
	v1Client "github.com/roadrunner-server/api-go/v6/centrifugo/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// stubAPI implements only Publish; every other CentrifugoApi method responds
// with codes.Unimplemented through the embedded unimplemented server, which
// lets the tests cover both forward() outcomes.
type stubAPI struct {
	v1Client.UnimplementedCentrifugoApiServer
}

func (stubAPI) Publish(context.Context, *v1Client.PublishRequest) (*v1Client.PublishResponse, error) {
	return &v1Client.PublishResponse{}, nil
}

func newBufconnRPC(t *testing.T) *rpc {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	v1Client.RegisterCentrifugoApiServer(srv, stubAPI{})

	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return &rpc{client: &client{centrifugoClient: v1Client.NewCentrifugoApiClient(conn)}, log: testLogger()}
}

func TestRPCForwardSuccess(t *testing.T) {
	r := newBufconnRPC(t)

	resp, err := r.Publish(t.Context(), connect.NewRequest(&v1Client.PublishRequest{}))
	require.NoError(t, err)
	require.NotNil(t, resp.Msg)
}

func TestRPCForwardMapsGRPCStatus(t *testing.T) {
	r := newBufconnRPC(t)

	calls := map[string]func(ctx context.Context) error{
		"batch": func(ctx context.Context) error {
			_, err := r.Batch(ctx, connect.NewRequest(&v1Client.BatchRequest{}))
			return err
		},
		"broadcast": func(ctx context.Context) error {
			_, err := r.Broadcast(ctx, connect.NewRequest(&v1Client.BroadcastRequest{}))
			return err
		},
		"subscribe": func(ctx context.Context) error {
			_, err := r.Subscribe(ctx, connect.NewRequest(&v1Client.SubscribeRequest{}))
			return err
		},
		"unsubscribe": func(ctx context.Context) error {
			_, err := r.Unsubscribe(ctx, connect.NewRequest(&v1Client.UnsubscribeRequest{}))
			return err
		},
		"disconnect": func(ctx context.Context) error {
			_, err := r.Disconnect(ctx, connect.NewRequest(&v1Client.DisconnectRequest{}))
			return err
		},
		"presence": func(ctx context.Context) error {
			_, err := r.Presence(ctx, connect.NewRequest(&v1Client.PresenceRequest{}))
			return err
		},
		"presence_stats": func(ctx context.Context) error {
			_, err := r.PresenceStats(ctx, connect.NewRequest(&v1Client.PresenceStatsRequest{}))
			return err
		},
		"history": func(ctx context.Context) error {
			_, err := r.History(ctx, connect.NewRequest(&v1Client.HistoryRequest{}))
			return err
		},
		"history_remove": func(ctx context.Context) error {
			_, err := r.HistoryRemove(ctx, connect.NewRequest(&v1Client.HistoryRemoveRequest{}))
			return err
		},
		"info": func(ctx context.Context) error {
			_, err := r.Info(ctx, connect.NewRequest(&v1Client.InfoRequest{}))
			return err
		},
		"rpc": func(ctx context.Context) error {
			_, err := r.RPC(ctx, connect.NewRequest(&v1Client.RPCRequest{}))
			return err
		},
		"refresh": func(ctx context.Context) error {
			_, err := r.Refresh(ctx, connect.NewRequest(&v1Client.RefreshRequest{}))
			return err
		},
		"channels": func(ctx context.Context) error {
			_, err := r.Channels(ctx, connect.NewRequest(&v1Client.ChannelsRequest{}))
			return err
		},
		"connections": func(ctx context.Context) error {
			_, err := r.Connections(ctx, connect.NewRequest(&v1Client.ConnectionsRequest{}))
			return err
		},
		"update_user_status": func(ctx context.Context) error {
			_, err := r.UpdateUserStatus(ctx, connect.NewRequest(&v1Client.UpdateUserStatusRequest{}))
			return err
		},
		"get_user_status": func(ctx context.Context) error {
			_, err := r.GetUserStatus(ctx, connect.NewRequest(&v1Client.GetUserStatusRequest{}))
			return err
		},
		"delete_user_status": func(ctx context.Context) error {
			_, err := r.DeleteUserStatus(ctx, connect.NewRequest(&v1Client.DeleteUserStatusRequest{}))
			return err
		},
		"block_user": func(ctx context.Context) error {
			_, err := r.BlockUser(ctx, connect.NewRequest(&v1Client.BlockUserRequest{}))
			return err
		},
		"unblock_user": func(ctx context.Context) error {
			_, err := r.UnblockUser(ctx, connect.NewRequest(&v1Client.UnblockUserRequest{}))
			return err
		},
		"revoke_token": func(ctx context.Context) error {
			_, err := r.RevokeToken(ctx, connect.NewRequest(&v1Client.RevokeTokenRequest{}))
			return err
		},
		"invalidate_user_tokens": func(ctx context.Context) error {
			_, err := r.InvalidateUserTokens(ctx, connect.NewRequest(&v1Client.InvalidateUserTokensRequest{}))
			return err
		},
		"device_register": func(ctx context.Context) error {
			_, err := r.DeviceRegister(ctx, connect.NewRequest(&v1Client.DeviceRegisterRequest{}))
			return err
		},
		"device_update": func(ctx context.Context) error {
			_, err := r.DeviceUpdate(ctx, connect.NewRequest(&v1Client.DeviceUpdateRequest{}))
			return err
		},
		"device_remove": func(ctx context.Context) error {
			_, err := r.DeviceRemove(ctx, connect.NewRequest(&v1Client.DeviceRemoveRequest{}))
			return err
		},
		"device_list": func(ctx context.Context) error {
			_, err := r.DeviceList(ctx, connect.NewRequest(&v1Client.DeviceListRequest{}))
			return err
		},
		"device_topic_list": func(ctx context.Context) error {
			_, err := r.DeviceTopicList(ctx, connect.NewRequest(&v1Client.DeviceTopicListRequest{}))
			return err
		},
		"device_topic_update": func(ctx context.Context) error {
			_, err := r.DeviceTopicUpdate(ctx, connect.NewRequest(&v1Client.DeviceTopicUpdateRequest{}))
			return err
		},
		"user_topic_list": func(ctx context.Context) error {
			_, err := r.UserTopicList(ctx, connect.NewRequest(&v1Client.UserTopicListRequest{}))
			return err
		},
		"user_topic_update": func(ctx context.Context) error {
			_, err := r.UserTopicUpdate(ctx, connect.NewRequest(&v1Client.UserTopicUpdateRequest{}))
			return err
		},
		"send_push_notification": func(ctx context.Context) error {
			_, err := r.SendPushNotification(ctx, connect.NewRequest(&v1Client.SendPushNotificationRequest{}))
			return err
		},
		"update_push_status": func(ctx context.Context) error {
			_, err := r.UpdatePushStatus(ctx, connect.NewRequest(&v1Client.UpdatePushStatusRequest{}))
			return err
		},
		"cancel_push": func(ctx context.Context) error {
			_, err := r.CancelPush(ctx, connect.NewRequest(&v1Client.CancelPushRequest{}))
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err := call(t.Context())
			require.Error(t, err)
			// the stub leaves the method unimplemented; the gRPC status code
			// must come back as the matching connect code
			assert.Equal(t, connect.CodeUnimplemented, connect.CodeOf(err))
		})
	}
}

func TestPluginRPCHandler(t *testing.T) {
	p := &Plugin{client: &client{}, log: testLogger()}

	path, handler := p.RPC()
	require.NotEmpty(t, path)
	require.NotNil(t, handler)
}
