//go:build linux || darwin || freebsd

package centrifugo

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"slices"
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
		wantErr string
	}{
		{name: "default TCP"},
		{name: "UNIX defaults", addr: "unix://test.sock"},
		{name: "empty options", addr: "unix://test.sock", options: "{}"},
		{name: "mode only", addr: "unix://test.sock", options: `{mode: "0600"}`},
		{name: "explicit zero", addr: "unix://test.sock", options: `{mode: "0000", uid: 0, gid: 0}`},
		{name: "unset mode", addr: "unix://test.sock", options: "{uid: 0, gid: 0}"},
		{name: "default TCP empty options", options: "{}"},
		{name: "default TCP options", options: `{mode: "0600"}`, wantErr: "filesystem unix:// address"},
		{name: "TCP options", addr: "tcp://127.0.0.1:0", options: `{mode: "0600"}`, wantErr: "filesystem unix:// address"},
		{name: "empty socket path", addr: "unix://", options: `{mode: "0600"}`, wantErr: "filesystem unix:// address"},
		{name: "invalid mode", addr: "unix://test.sock", options: `{mode: "0780"}`, wantErr: "invalid unix socket mode"},
		{name: "unquoted mode", addr: "unix://test.sock", options: "{mode: 0660}", wantErr: "invalid unix socket mode"},
		{name: "scalar options", addr: "unix://test.sock", options: "false", wantErr: "expected a map"},
		{name: "negative UID", addr: "unix://test.sock", options: "{uid: -1}", wantErr: "invalid unix socket uid"},
		{name: "negative GID", addr: "unix://test.sock", options: "{gid: -1}", wantErr: "invalid unix socket gid"},
		{name: "reserved UID", addr: "unix://test.sock", options: "{uid: 4294967295}", wantErr: "invalid unix socket uid"},
		{name: "reserved GID", addr: "unix://test.sock", options: "{gid: 4294967295}", wantErr: "invalid unix socket gid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".rr.yaml")
			data := fmt.Sprintf(`version: "3"
centrifuge:
  proxy_address: %q
`, tc.addr)
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
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
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

func TestProxySocketOwnershipError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Requires an unprivileged process.")
	}
	worker, err := filepath.Abs("php_test_files/centrifuge_connect.php")
	require.NoError(t, err)
	groups, err := os.Getgroups()
	require.NoError(t, err)
	otherGID := 0
	for otherGID == os.Getegid() || slices.Contains(groups, otherGID) {
		otherGID++
	}

	for _, tc := range []struct {
		field string
		id    int
	}{
		{field: "uid", id: 0},
		{field: "gid", id: otherGID},
	} {
		t.Run(tc.field, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("RR_TEST_SOCKET_ID", strconv.Itoa(tc.id))
			data := fmt.Sprintf(`version: "3"
server:
  command: [php, %q]
centrifuge:
  proxy_address: unix://ownership.sock
  proxy_socket: {%s: "${RR_TEST_SOCKET_ID}"}
  pool:
    num_workers: 1
    destroy_timeout: 5s
`, worker, tc.field)
			require.NoError(t, os.WriteFile(".rr.yaml", []byte(data), 0o600))
			cfg := &config.Plugin{Path: ".rr.yaml"}
			require.NoError(t, cfg.Init())
			log := &logger.Plugin{}
			require.NoError(t, log.Init(cfg))
			rrServer := &server.Plugin{}
			require.NoError(t, rrServer.Init(cfg, log.ServiceLogger()))
			t.Cleanup(func() { require.NoError(t, rrServer.Stop(context.Background())) })
			p := &centrifuge.Plugin{}
			require.NoError(t, p.Init(cfg, log.ServiceLogger(), rrServer))
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				require.NoError(t, p.Stop(ctx))
			})
			select {
			case errS := <-p.Serve():
				require.ErrorContains(t, errS, "chown unix socket")
			case <-time.After(5 * time.Second):
				t.Fatal("Centrifuge did not report the ownership error")
			}
			_, errS := os.Stat("ownership.sock")
			require.ErrorIs(t, errS, os.ErrNotExist)
		})
	}
}
