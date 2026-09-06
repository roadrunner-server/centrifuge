//go:build linux || darwin || freebsd

package centrifugo

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	centrifugov1 "github.com/roadrunner-server/api-go/v6/centrifugo/proxy/v1"
	"github.com/roadrunner-server/centrifuge/v6"
	"github.com/roadrunner-server/config/v6"
	"github.com/roadrunner-server/endure/v2"
	"github.com/roadrunner-server/logger/v6"
	"github.com/roadrunner-server/server/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestProxySocketConfig(t *testing.T) {
	for _, tc := range []struct {
		name    string
		addr    string
		options string
		mode    string
		zeroIDs bool
		invalid bool
	}{
		{name: "default TCP"},
		{name: "UNIX defaults", addr: "unix://test.sock"},
		{name: "empty options", addr: "unix://test.sock", options: "{}"},
		{name: "mode only", addr: "unix://test.sock", options: `{mode: "0600"}`, mode: "0600"},
		{name: "explicit zero", addr: "unix://test.sock", options: `{mode: "0000", uid: 0, gid: 0}`, mode: "0000", zeroIDs: true},
		{name: "unset mode", addr: "unix://test.sock", options: "{uid: 0, gid: 0}", zeroIDs: true},
		{name: "default TCP options", options: "{}", invalid: true},
		{name: "unquoted mode", addr: "unix://test.sock", options: "{mode: 0660}", invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".rr.yaml")
			data := fmt.Sprintf("version: '3'\ncentrifuge:\n  proxy_address: %q\n", tc.addr)
			if tc.options != "" {
				data += "  proxy_socket: " + tc.options + "\n"
			}
			require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
			cfg := &config.Plugin{Path: path}
			require.NoError(t, cfg.Init())
			log := &logger.Plugin{}
			require.NoError(t, log.Init(cfg))
			p := &centrifuge.Plugin{}
			err := p.Init(cfg, log.ServiceLogger(), nil)
			if tc.invalid {
				require.ErrorContains(t, err, "centrifuge.proxy_socket")
				return
			}
			require.NoError(t, err)
			var decoded centrifuge.Config
			require.NoError(t, cfg.UnmarshalKey("centrifuge", &decoded))
			if tc.options == "{}" {
				require.True(t, cfg.Has("centrifuge.proxy_socket"))
				require.Nil(t, decoded.ProxySocket)
				return
			}
			require.NoError(t, decoded.InitDefaults())
			require.Equal(t, "127.0.0.1:10000", decoded.GrpcAPIAddress)
			options := decoded.ProxySocket
			if tc.options == "" {
				require.Nil(t, options)
				return
			}
			require.NotNil(t, options)
			require.Equal(t, tc.mode, options.Mode)
			if tc.zeroIDs {
				require.NotNil(t, options.UID)
				require.NotNil(t, options.GID)
				require.Zero(t, *options.UID)
				require.Zero(t, *options.GID)
			} else {
				require.Nil(t, options.UID)
				require.Nil(t, options.GID)
			}
		})
	}
}

func TestProxySocketRawIDs(t *testing.T) {
	for _, field := range []string{"uid", "gid"} {
		for _, tc := range []struct {
			value   string
			env     string
			want    int
			invalid bool
		}{
			{value: "null"},
			{value: "0"},
			{value: "33.0", want: 33},
			{value: `"0x21"`, want: 33},
			{value: `"${RR_TEST_SOCKET_ID}"`, env: "0"},
			{value: `"${RR_TEST_SOCKET_ID}"`, env: "33", want: 33},
			{value: `"${RR_TEST_SOCKET_ID}"`, invalid: true},
			{value: `""`, invalid: true},
			{value: "1.9", invalid: true},
			{value: "-0.5", invalid: true},
			{value: "true", invalid: true},
			{value: "false", invalid: true},
			{value: "-1", invalid: true},
			{value: "4294967295", invalid: true},
			{value: `"4294967295"`, invalid: true},
			{value: "[]", invalid: true},
			{value: "{}", invalid: true},
		} {
			t.Run(field+"/"+tc.value+"/"+tc.env, func(t *testing.T) {
				t.Setenv("RR_TEST_SOCKET_ID", tc.env)
				if tc.env == "" {
					require.NoError(t, os.Unsetenv("RR_TEST_SOCKET_ID"))
				}
				dir := t.TempDir()
				socket := filepath.Join(dir, "proxy.sock")
				path := filepath.Join(dir, ".rr.json")
				data := fmt.Sprintf(`{"version":"3","centrifuge":{
					"proxy_address":%q,"proxy_socket":{%q:%s},
					"grpc_api_address":"tcp://127.0.0.1:10000"}}`, "unix://"+socket, field, tc.value)
				require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
				cfg := &config.Plugin{Path: path}
				require.NoError(t, cfg.Init())
				log := &logger.Plugin{}
				require.NoError(t, log.Init(cfg))
				p := &centrifuge.Plugin{}
				err := p.Init(cfg, log.ServiceLogger(), nil)
				_, statErr := os.Stat(socket)
				require.ErrorIs(t, statErr, os.ErrNotExist)
				if tc.invalid {
					require.ErrorContains(t, err, "centrifuge.proxy_socket."+field)
					return
				}
				require.NoError(t, err)
				var decoded centrifuge.Config
				require.NoError(t, cfg.UnmarshalKey("centrifuge", &decoded))
				require.NoError(t, decoded.InitDefaults())
				require.Equal(t, "127.0.0.1:10000", decoded.GrpcAPIAddress)
				options := decoded.ProxySocket
				if tc.value == "null" {
					if options != nil {
						require.Nil(t, options.UID)
						require.Nil(t, options.GID)
					}
					return
				}
				require.NotNil(t, options)
				id, unset := options.UID, options.GID
				if field == "gid" {
					id, unset = options.GID, options.UID
				}
				require.NotNil(t, id)
				require.Equal(t, tc.want, *id)
				require.Nil(t, unset)
			})
		}
	}
}

func TestProxyUnixSocket(t *testing.T) {
	worker, err := filepath.Abs("php_test_files/centrifuge_connect.php")
	require.NoError(t, err)
	t.Chdir(t.TempDir())
	uid, gid := os.Getuid(), os.Getgid()
	t.Setenv("RR_TEST_SOCKET_UID", strconv.Itoa(uid))
	t.Setenv("RR_TEST_SOCKET_GID", strconv.Itoa(gid))
	data := fmt.Sprintf(`
version: "3"
server:
  command: [php, %q]
centrifuge:
  proxy_address: unix://proxy.sock
  proxy_socket: {mode: "0640", uid: "${RR_TEST_SOCKET_UID}", gid: "${RR_TEST_SOCKET_GID}"}
  grpc_api_address: tcp://127.0.0.1:10000
  pool:
    num_workers: 1
    destroy_timeout: 5s
`, worker)
	require.NoError(t, os.WriteFile(".rr.yaml", []byte(data), 0o600))
	cfg := &config.Plugin{Path: ".rr.yaml"}
	cont := endure.New(slog.LevelError)
	require.NoError(t, cont.RegisterAll(cfg, &logger.Plugin{}, &server.Plugin{}, &centrifuge.Plugin{}))
	require.NoError(t, cont.Init())
	errCh, err := cont.Serve()
	require.NoError(t, err)
	stop := sync.OnceValue(cont.Stop)
	t.Cleanup(func() { require.NoError(t, stop()) })

	// The incoming proxy does not need an outbound Centrifugo API call.
	conn, err := grpc.NewClient("passthrough:///proxy", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", "proxy.sock")
		}))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	response, err := centrifugov1.NewCentrifugoProxyClient(conn).Connect(ctx, &centrifugov1.ConnectRequest{Client: "unix-test"}, grpc.WaitForReady(true))
	require.NoError(t, err)
	require.Nil(t, response.GetError())
	require.Equal(t, "ca94a2db-b5e8-4551-986b-5bad3e871fd1", response.GetResult().GetUser())

	info, err := os.Stat("proxy.sock")
	require.NoError(t, err)
	require.NotZero(t, info.Mode()&os.ModeSocket)
	require.Equal(t, os.FileMode(0o640), info.Mode().Perm())
	stat := info.Sys().(*syscall.Stat_t)
	require.EqualValues(t, uid, stat.Uid)
	require.EqualValues(t, gid, stat.Gid)
	require.NoError(t, conn.Close())
	select {
	case result := <-errCh:
		t.Fatalf("serve error: %v", result)
	default:
	}
	require.NoError(t, stop())
	_, err = os.Stat("proxy.sock")
	require.ErrorIs(t, err, os.ErrNotExist)
}
