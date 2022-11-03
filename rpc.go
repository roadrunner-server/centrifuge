package centrifuge

import (
	"context"

	v1Client "go.buf.build/grpc/go/roadrunner-server/api/centrifugo/api/v1"
	"go.uber.org/zap"
)

type rpc struct {
	client *client
	log    *zap.Logger
}

/*
service CentrifugoApi {
    rpc Publish (PublishRequest) returns (PublishResponse) {}
    rpc Broadcast (BroadcastRequest) returns (BroadcastResponse) {}
    rpc Subscribe (SubscribeRequest) returns (SubscribeResponse) {}
    rpc Unsubscribe (UnsubscribeRequest) returns (UnsubscribeResponse) {}
    rpc Disconnect (DisconnectRequest) returns (DisconnectResponse) {}
    rpc Presence (PresenceRequest) returns (PresenceResponse) {}
    rpc PresenceStats (PresenceStatsRequest) returns (PresenceStatsResponse) {}
    rpc History (HistoryRequest) returns (HistoryResponse) {}
    rpc HistoryRemove (HistoryRemoveRequest) returns (HistoryRemoveResponse) {}
    rpc Info (InfoRequest) returns (InfoResponse) {}
    rpc RPC (RPCRequest) returns (RPCResponse) {}
    rpc Refresh (RefreshRequest) returns (RefreshResponse) {}
    rpc Channels (ChannelsRequest) returns (ChannelsResponse) {}
    rpc Connections (ConnectionsRequest) returns (ConnectionsResponse) {}
    rpc UpdateUserStatus (UpdateUserStatusRequest) returns (UpdateUserStatusResponse) {}
    rpc GetUserStatus (GetUserStatusRequest) returns (GetUserStatusResponse) {}
    rpc DeleteUserStatus (DeleteUserStatusRequest) returns (DeleteUserStatusResponse) {}
    rpc BlockUser (BlockUserRequest) returns (BlockUserResponse) {}
    rpc UnblockUser (UnblockUserRequest) returns (UnblockUserResponse) {}
    rpc RevokeToken (RevokeTokenRequest) returns (RevokeTokenResponse) {}
    rpc InvalidateUserTokens (InvalidateUserTokensRequest) returns (InvalidateUserTokensResponse) {}
}
*/

func (r *rpc) Publish(in *v1Client.PublishRequest, out *v1Client.PublishResponse) error {
	r.log.Debug("got publish request")
	resp, err := r.client.client().Publish(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Broadcast(in *v1Client.BroadcastRequest, out *v1Client.BroadcastResponse) error {
	r.log.Debug("got broadcast request")
	resp, err := r.client.client().Broadcast(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Subscribe(in *v1Client.SubscribeRequest, out *v1Client.SubscribeResponse) error {
	r.log.Debug("got subscribe request")
	resp, err := r.client.client().Subscribe(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Unsubscribe(in *v1Client.UnsubscribeRequest, out *v1Client.UnsubscribeResponse) error {
	r.log.Debug("got unsubscribe request")
	resp, err := r.client.client().Unsubscribe(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Disconnect(in *v1Client.DisconnectRequest, out *v1Client.DisconnectResponse) error {
	r.log.Debug("got disconnect request")
	resp, err := r.client.client().Disconnect(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}
func (r *rpc) Presence(in *v1Client.PresenceRequest, out *v1Client.PresenceResponse) error {
	r.log.Debug("got presence request")
	resp, err := r.client.client().Presence(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) PresenceStats(in *v1Client.PresenceStatsRequest, out *v1Client.PresenceStatsResponse) error {
	r.log.Debug("got presence_stats request")
	resp, err := r.client.client().PresenceStats(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) History(in *v1Client.HistoryRequest, out *v1Client.HistoryResponse) error {
	r.log.Debug("got history request")
	resp, err := r.client.client().History(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) HistoryRemove(in *v1Client.HistoryRemoveRequest, out *v1Client.HistoryRemoveResponse) error {
	r.log.Debug("got history_remove request")
	resp, err := r.client.client().HistoryRemove(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Info(in *v1Client.InfoRequest, out *v1Client.InfoResponse) error {
	r.log.Debug("got info request")
	resp, err := r.client.client().Info(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Refresh(in *v1Client.RefreshRequest, out *v1Client.RefreshResponse) error {
	r.log.Debug("got refresh request")
	resp, err := r.client.client().Refresh(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Channels(in *v1Client.ChannelsRequest, out *v1Client.ChannelsResponse) error {
	r.log.Debug("got channels request")
	resp, err := r.client.client().Channels(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) Connections(in *v1Client.ConnectionsRequest, out *v1Client.ConnectionsResponse) error {
	r.log.Debug("got connections request")
	resp, err := r.client.client().Connections(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UpdateUserStatus(in *v1Client.UpdateUserStatusRequest, out *v1Client.UpdateUserStatusResponse) error {
	r.log.Debug("got update_user_status request")
	resp, err := r.client.client().UpdateUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) GetUserStatus(in *v1Client.GetUserStatusRequest, out *v1Client.GetUserStatusResponse) error {
	r.log.Debug("got get_user_status request")
	resp, err := r.client.client().GetUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) DeleteUserStatus(in *v1Client.DeleteUserStatusRequest, out *v1Client.DeleteUserStatusResponse) error {
	r.log.Debug("got delete_user_status request")
	resp, err := r.client.client().DeleteUserStatus(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) BlockUser(in *v1Client.BlockUserRequest, out *v1Client.BlockUserResponse) error {
	r.log.Debug("got block_user request")
	resp, err := r.client.client().BlockUser(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) UnblockUser(in *v1Client.UnblockUserRequest, out *v1Client.UnblockUserResponse) error {
	r.log.Debug("got unblock_user request")
	resp, err := r.client.client().UnblockUser(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) RevokeToken(in *v1Client.RevokeTokenRequest, out *v1Client.RevokeTokenResponse) error {
	r.log.Debug("got revoke_token request")
	resp, err := r.client.client().RevokeToken(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}

func (r *rpc) InvalidateUserTokens(in *v1Client.InvalidateUserTokensRequest, out *v1Client.InvalidateUserTokensResponse) error {
	r.log.Debug("got invalidate_user_tokens request")
	resp, err := r.client.client().InvalidateUserTokens(context.Background(), in)
	if err != nil {
		return err
	}

	out.Error = resp.GetError()
	out.Result = resp.GetResult()

	return nil
}
