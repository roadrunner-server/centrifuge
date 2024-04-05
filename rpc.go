package centrifuge

import (
	"context"

	v1Client "github.com/roadrunner-server/api/v4/build/centrifugo/api/v1"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
)

type rpc struct {
	client *client
	log    *zap.Logger
}

/*
service CentrifugoApi {
  rpc Batch(BatchRequest) returns (BatchResponse) {}
  rpc Publish(PublishRequest) returns (PublishResponse) {}
  rpc Broadcast(BroadcastRequest) returns (BroadcastResponse) {}
  rpc Subscribe(SubscribeRequest) returns (SubscribeResponse) {}
  rpc Unsubscribe(UnsubscribeRequest) returns (UnsubscribeResponse) {}
  rpc Disconnect(DisconnectRequest) returns (DisconnectResponse) {}
  rpc Presence(PresenceRequest) returns (PresenceResponse) {}
  rpc PresenceStats(PresenceStatsRequest) returns (PresenceStatsResponse) {}
  rpc History(HistoryRequest) returns (HistoryResponse) {}
  rpc HistoryRemove(HistoryRemoveRequest) returns (HistoryRemoveResponse) {}
  rpc Info(InfoRequest) returns (InfoResponse) {}
  rpc RPC(RPCRequest) returns (RPCResponse) {}
  rpc Refresh(RefreshRequest) returns (RefreshResponse) {}
  rpc Channels(ChannelsRequest) returns (ChannelsResponse) {}
  rpc Connections(ConnectionsRequest) returns (ConnectionsResponse) {}
  rpc UpdateUserStatus(UpdateUserStatusRequest) returns (UpdateUserStatusResponse) {}
  rpc GetUserStatus(GetUserStatusRequest) returns (GetUserStatusResponse) {}
  rpc DeleteUserStatus(DeleteUserStatusRequest) returns (DeleteUserStatusResponse) {}
  rpc BlockUser(BlockUserRequest) returns (BlockUserResponse) {}
  rpc UnblockUser(UnblockUserRequest) returns (UnblockUserResponse) {}
  rpc RevokeToken(RevokeTokenRequest) returns (RevokeTokenResponse) {}
  rpc InvalidateUserTokens(InvalidateUserTokensRequest) returns (InvalidateUserTokensResponse) {}
  rpc DeviceRegister(DeviceRegisterRequest) returns (DeviceRegisterResponse) {}
  rpc DeviceUpdate(DeviceUpdateRequest) returns (DeviceUpdateResponse) {}
  rpc DeviceRemove(DeviceRemoveRequest) returns (DeviceRemoveResponse) {}
  rpc DeviceList(DeviceListRequest) returns (DeviceListResponse) {}
  rpc DeviceTopicList(DeviceTopicListRequest) returns (DeviceTopicListResponse) {}
  rpc DeviceTopicUpdate(DeviceTopicUpdateRequest) returns (DeviceTopicUpdateResponse) {}
  rpc UserTopicList(UserTopicListRequest) returns (UserTopicListResponse) {}
  rpc UserTopicUpdate(UserTopicUpdateRequest) returns (UserTopicUpdateResponse) {}
  rpc SendPushNotification(SendPushNotificationRequest) returns (SendPushNotificationResponse) {}
  rpc UpdatePushStatus(UpdatePushStatusRequest) returns (UpdatePushStatusResponse) {}
  rpc CancelPush(CancelPushRequest) returns (CancelPushResponse) {}
  rpc RateLimit(RateLimitRequest) returns (RateLimitResponse) {}
}
*/

func (r *rpc) Batch(in *v1Client.BatchRequest, out *v1Client.BatchResponse) error {
	r.log.Debug("got butch request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}

	resp, err := client.Batch(context.Background(), in)
	if err != nil {
		return err
	}

	out.Replies = resp.GetReplies()

	return nil
}

func (r *rpc) Publish(in *v1Client.PublishRequest, out *v1Client.PublishResponse) error {
	r.log.Debug("got publish request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}

	resp, err := client.Publish(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Broadcast(in *v1Client.BroadcastRequest, out *v1Client.BroadcastResponse) error {
	r.log.Debug("got broadcast request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Broadcast(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Subscribe(in *v1Client.SubscribeRequest, out *v1Client.SubscribeResponse) error {
	r.log.Debug("got subscribe request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Subscribe(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Unsubscribe(in *v1Client.UnsubscribeRequest, out *v1Client.UnsubscribeResponse) error {
	r.log.Debug("got unsubscribe request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Unsubscribe(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Disconnect(in *v1Client.DisconnectRequest, out *v1Client.DisconnectResponse) error {
	r.log.Debug("got disconnect request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Disconnect(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}
func (r *rpc) Presence(in *v1Client.PresenceRequest, out *v1Client.PresenceResponse) error {
	r.log.Debug("got presence request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Presence(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) PresenceStats(in *v1Client.PresenceStatsRequest, out *v1Client.PresenceStatsResponse) error {
	r.log.Debug("got presence_stats request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.PresenceStats(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) History(in *v1Client.HistoryRequest, out *v1Client.HistoryResponse) error {
	r.log.Debug("got history request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.History(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) HistoryRemove(in *v1Client.HistoryRemoveRequest, out *v1Client.HistoryRemoveResponse) error {
	r.log.Debug("got history_remove request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.HistoryRemove(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Info(in *v1Client.InfoRequest, out *v1Client.InfoResponse) error {
	r.log.Debug("got info request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Info(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) RPC(in *v1Client.RPCRequest, out *v1Client.RPCResponse) error {
	r.log.Debug("got rpc request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.RPC(context.Background(), in)
	if err != nil {
		return err
	}

	out.Result = resp.GetResult()
	out.Error = resp.GetError()

	return nil
}

func (r *rpc) Refresh(in *v1Client.RefreshRequest, out *v1Client.RefreshResponse) error {
	r.log.Debug("got refresh request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Refresh(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Channels(in *v1Client.ChannelsRequest, out *v1Client.ChannelsResponse) error {
	r.log.Debug("got channels request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Channels(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Connections(in *v1Client.ConnectionsRequest, out *v1Client.ConnectionsResponse) error {
	r.log.Debug("got connections request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.Connections(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UpdateUserStatus(in *v1Client.UpdateUserStatusRequest, out *v1Client.UpdateUserStatusResponse) error {
	r.log.Debug("got update_user_status request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.UpdateUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) GetUserStatus(in *v1Client.GetUserStatusRequest, out *v1Client.GetUserStatusResponse) error {
	r.log.Debug("got get_user_status request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.GetUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeleteUserStatus(in *v1Client.DeleteUserStatusRequest, out *v1Client.DeleteUserStatusResponse) error {
	r.log.Debug("got delete_user_status request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeleteUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) BlockUser(in *v1Client.BlockUserRequest, out *v1Client.BlockUserResponse) error {
	r.log.Debug("got block_user request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.BlockUser(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UnblockUser(in *v1Client.UnblockUserRequest, out *v1Client.UnblockUserResponse) error {
	r.log.Debug("got unblock_user request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.UnblockUser(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) RevokeToken(in *v1Client.RevokeTokenRequest, out *v1Client.RevokeTokenResponse) error {
	r.log.Debug("got revoke_token request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.RevokeToken(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) InvalidateUserTokens(in *v1Client.InvalidateUserTokensRequest, out *v1Client.InvalidateUserTokensResponse) error {
	r.log.Debug("got invalidate_user_tokens request")
	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.InvalidateUserTokens(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceRegister(in *v1Client.DeviceRegisterRequest, out *v1Client.DeviceRegisterResponse) error {
	r.log.Debug("got device register request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceRegister(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceUpdate(in *v1Client.DeviceUpdateRequest, out *v1Client.DeviceUpdateResponse) error {
	r.log.Debug("got device update request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceUpdate(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceRemove(in *v1Client.DeviceRemoveRequest, out *v1Client.DeviceRemoveResponse) error {
	r.log.Debug("got device remove request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceRemove(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceList(in *v1Client.DeviceListRequest, out *v1Client.DeviceListResponse) error {
	r.log.Debug("got device remove request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceList(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceTopicList(in *v1Client.DeviceTopicListRequest, out *v1Client.DeviceTopicListResponse) error {
	r.log.Debug("got device topic list request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceTopicList(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeviceTopicUpdate(in *v1Client.DeviceTopicUpdateRequest, out *v1Client.DeviceTopicUpdateResponse) error {
	r.log.Debug("got device topic update request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.DeviceTopicUpdate(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UserTopicList(in *v1Client.UserTopicListRequest, out *v1Client.UserTopicListResponse) error {
	r.log.Debug("got user topic list request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.UserTopicList(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UserTopicUpdate(in *v1Client.UserTopicUpdateRequest, out *v1Client.UserTopicUpdateResponse) error {
	r.log.Debug("got user topic update request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.UserTopicUpdate(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) SendPushNotification(in *v1Client.SendPushNotificationRequest, out *v1Client.SendPushNotificationResponse) error {
	r.log.Debug("got send push notification request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.SendPushNotification(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UpdatePushStatus(in *v1Client.UpdatePushStatusRequest, out *v1Client.UpdatePushStatusResponse) error {
	r.log.Debug("got update push status request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.UpdatePushStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) CancelPush(in *v1Client.CancelPushRequest, out *v1Client.CancelPushResponse) error {
	r.log.Debug("got cancel push request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.CancelPush(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) RateLimit(in *v1Client.RateLimitRequest, out *v1Client.RateLimitResponse) error {
	r.log.Debug("got rate limit request")

	client := r.client.client()
	if client == nil {
		return errors.Str("RoadRunner is not ready yet, try in a few seconds")
	}
	resp, err := client.RateLimit(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}
