package centrifugo

import (
	"testing"
	"time"

	"tests/helpers"

	centrifugeClient "github.com/centrifugal/centrifuge-go"
	"github.com/roadrunner-server/centrifuge/v6"
	rpcPlugin "github.com/roadrunner-server/rpc/v6"
	"github.com/roadrunner-server/server/v6"
	"github.com/stretchr/testify/require"
)

const (
	// centrifugoWS is the websocket endpoint the broker serves; the proxy
	// endpoints in env/config.json point back at the plugin on 10001.
	centrifugoWS   = "127.0.0.1:8000"
	proxyAddr      = "127.0.0.1:10001"
	clientTimeout  = time.Second * 100
	logWait        = time.Second * 15
	logWaitTick    = time.Millisecond * 50
	subscribeTopic = "test"
)

func centrifugePlugins() []any {
	return []any{&server.Plugin{}, &rpcPlugin.Plugin{}, &centrifuge.Plugin{}}
}

// TestProxiesConnectAndSubscribe drives a real centrifugo client through the
// broker, which proxies both calls back to the plugin over grpc. The assertions
// are the proxy records the plugin logs, so a broken proxy path fails rather
// than merely logging nothing.
func TestProxiesConnectAndSubscribe(t *testing.T) {
	helpers.WaitForCentrifugo(t, centrifugoWS)

	rr, _ := helpers.Start(t,
		"configs/.rr-centrifugo-init.yaml",
		centrifugePlugins(),
		helpers.WithObservedLogger(),
		helpers.WithTCPProbe(proxyAddr),
	)

	client := centrifugeClient.NewProtobufClient("ws://"+centrifugoWS+"/connection/websocket", centrifugeClient.Config{
		Name:               "roadrunner_tests",
		Version:            "3.0.0",
		ReadTimeout:        clientTimeout,
		WriteTimeout:       clientTimeout,
		HandshakeTimeout:   clientTimeout,
		MaxServerPingDelay: clientTimeout,
	})
	t.Cleanup(client.Close)

	require.NoError(t, client.Connect())

	require.Eventually(t, func() bool {
		return rr.Logs.FilterMessageSnippet("got connect proxy request").Len() == 1
	}, logWait, logWaitTick, "the connect request never reached the plugin")

	subscription, err := client.NewSubscription(subscribeTopic)
	require.NoError(t, err)
	require.NoError(t, subscription.Subscribe())

	require.Eventually(t, func() bool {
		return rr.Logs.FilterMessageSnippet("got subscribe proxy request").Len() == 1
	}, logWait, logWaitTick, "the subscribe request never reached the plugin")

	require.NoError(t, subscription.Unsubscribe())

	// exactly one of each: a retrying proxy would show up as a higher count
	require.Equal(t, 1, rr.Logs.FilterMessageSnippet("got connect proxy request").Len())
	require.Equal(t, 1, rr.Logs.FilterMessageSnippet("got subscribe proxy request").Len())
}
