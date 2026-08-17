package centrifugo

import (
	"io"
	"net/http"
	"testing"

	"tests/helpers"

	"github.com/roadrunner-server/status/v6"
	"github.com/stretchr/testify/require"
)

const statusAddr = "127.0.0.1:35544"

// probeStatus fetches one of the status endpoints and returns its code and body.
func probeStatus(t *testing.T, path string) (int, string) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+statusAddr+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() { require.NoError(t, resp.Body.Close()) }()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(body)
}

// TestStatusEndpoints exercises every status assertion against a single running
// container. Booting one container per case raced the previous one for the
// status port, since the listener is not always released by the time the next
// container binds it.
func TestStatusEndpoints(t *testing.T) {
	helpers.Start(t,
		"configs/.rr-centrifugo-status.yaml",
		append(centrifugePlugins(), &status.Plugin{}),
		helpers.WithTCPProbe(statusAddr),
	)

	const healthy = `[{"plugin_name":"centrifuge","error_message":"","status_code":200}]`

	t.Run("health reports the pool is up", func(t *testing.T) {
		code, body := probeStatus(t, "/health?plugin=centrifuge")

		require.Equal(t, http.StatusOK, code)
		require.JSONEq(t, healthy, body)
	})

	t.Run("ready reports the pool is up", func(t *testing.T) {
		code, body := probeStatus(t, "/ready?plugin=centrifuge")

		require.Equal(t, http.StatusOK, code)
		require.JSONEq(t, healthy, body)
	})

	// an unregistered name must yield an empty list rather than centrifuge's
	// own entry, so the endpoint is not matching whatever it is asked for
	t.Run("unknown plugin reports nothing", func(t *testing.T) {
		code, body := probeStatus(t, "/health?plugin=not-registered")

		require.Equal(t, http.StatusOK, code)
		require.JSONEq(t, `[]`, body)
	})
}
