package centrifugo

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	mocklogger "tests/mock"

	centrifugeClient "github.com/centrifugal/centrifuge-go"
	"github.com/roadrunner-server/centrifuge/v5"
	"github.com/roadrunner-server/config/v5"
	"github.com/roadrunner-server/endure/v2"
	"github.com/roadrunner-server/logger/v5"
	"github.com/roadrunner-server/rpc/v5"
	"github.com/roadrunner-server/server/v5"
	"github.com/roadrunner-server/status/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCentrifugoPluginInit(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-centrifugo-init.yaml",
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("env/centrifugo_mac", "--config", "env/config.json", "--admin")
	} else {
		cmd = exec.Command("env/centrifugo", "--config", "env/config.json", "--admin")
	}

	err := cmd.Start()
	require.NoError(t, err)

	go func() {
		_ = cmd.Wait()
	}()

	l, oLogger := mocklogger.ZapTestLogger(zap.DebugLevel)
	err = cont.RegisterAll(
		l,
		cfg,
		&centrifuge.Plugin{},
		&server.Plugin{},
		&rpc.Plugin{},
	)
	assert.NoError(t, err)

	time.Sleep(time.Second)
	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := cont.Serve()
	if err != nil {
		t.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	stopCh := make(chan struct{}, 1)

	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-ch:
				assert.Fail(t, "error", e.Error.Error())
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}
			case <-sig:
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}

				return
			case <-stopCh:
				// timeout
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}

				return
			}
		}
	}()

	time.Sleep(time.Second * 5)
	client := centrifugeClient.NewProtobufClient("ws://127.0.0.1:8000/connection/websocket", centrifugeClient.Config{
		Name:               "roadrunner_tests",
		Version:            "3.0.0",
		ReadTimeout:        time.Second * 100,
		WriteTimeout:       time.Second * 100,
		HandshakeTimeout:   time.Second * 100,
		MaxServerPingDelay: time.Second * 100,
	})

	err = client.Connect()
	assert.NoError(t, err)

	subscription, err := client.NewSubscription("test")
	assert.NoError(t, err)

	err = subscription.Subscribe()
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)
	err = subscription.Unsubscribe()
	assert.NoError(t, err)
	client.Close()

	time.Sleep(time.Second * 2)
	stopCh <- struct{}{}
	_ = cmd.Process.Kill()
	wg.Wait()

	require.Equal(t, 1, oLogger.FilterMessageSnippet("got connect proxy request").Len())
	require.Equal(t, 1, oLogger.FilterMessageSnippet("got subscribe proxy request").Len())
}

func TestCentrifugoStatusChecks(t *testing.T) {
	cont := endure.New(slog.LevelDebug)

	cfg := &config.Plugin{
		Version: "2023.3.0",
		Path:    "configs/.rr-centrifugo-status.yaml",
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("env/centrifugo_mac", "--config", "env/config.json", "--admin")
	} else {
		cmd = exec.Command("env/centrifugo", "--config", "env/config.json", "--admin")
	}

	err := cmd.Start()
	require.NoError(t, err)

	go func() {
		_ = cmd.Wait()
	}()

	err = cont.RegisterAll(
		cfg,
		&logger.Plugin{},
		&status.Plugin{},
		&centrifuge.Plugin{},
		&server.Plugin{},
		&rpc.Plugin{},
	)
	assert.NoError(t, err)

	time.Sleep(time.Second)
	err = cont.Init()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := cont.Serve()
	if err != nil {
		t.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	stopCh := make(chan struct{}, 1)

	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-ch:
				assert.Fail(t, "error", e.Error.Error())
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}
			case <-sig:
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}

				return
			case <-stopCh:
				// timeout
				err = cont.Stop()
				if err != nil {
					assert.FailNow(t, "error", err.Error())
				}

				return
			}
		}
	}()

	time.Sleep(time.Second * 2)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://127.0.0.1:35544/health?plugin=centrifuge", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `[{"plugin_name":"centrifuge","error_message":"","status_code":200}]`, string(body))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "http://127.0.0.1:35544/ready?plugin=centrifuge", nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	body, _ = io.ReadAll(resp.Body)
	assert.Equal(t, `[{"plugin_name":"centrifuge","error_message":"","status_code":200}]`, string(body))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	time.Sleep(time.Second)
	stopCh <- struct{}{}
	wg.Wait()
}
