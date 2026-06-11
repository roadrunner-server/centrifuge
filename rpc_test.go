package centrifuge

import (
	"testing"

	"connectrpc.com/connect"
	v1Client "github.com/roadrunner-server/api-go/v6/centrifugo/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCBatchNotReady(t *testing.T) {
	// client() returns nil until connect() succeeds, so Batch returns early.
	r := &rpc{client: &client{}, log: testLogger()}

	_, err := r.Batch(t.Context(), connect.NewRequest(&v1Client.BatchRequest{}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RoadRunner is not ready yet")
	assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
}
