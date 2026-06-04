package centrifuge

import (
	"testing"

	v1Client "github.com/roadrunner-server/api-go/v6/centrifugo/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCBatchNotReady(t *testing.T) {
	// client() returns nil until connect() succeeds, so Batch returns early.
	r := &rpc{client: &client{}, log: testLogger()}

	err := r.Batch(&v1Client.BatchRequest{}, &v1Client.BatchResponse{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RoadRunner is not ready yet")
}
