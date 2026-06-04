package centrifuge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrpcConfig(t *testing.T) {
	cfg := &Config{GrpcAPIAddress: "tcp://foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "foo.bar", cfg.GrpcAPIAddress)
}

func TestGrpcConfig1(t *testing.T) {
	cfg := &Config{GrpcAPIAddress: "foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "foo.bar", cfg.GrpcAPIAddress)
}

func TestGrpcConfig2(t *testing.T) {
	cfg := &Config{GrpcAPIAddress: "tcp:/foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "tcp:/foo.bar", cfg.GrpcAPIAddress)
}

// TestGrpcConfigShortPrefix guards the off-by-one fixed by strings.CutPrefix:
// "tcp://x" is exactly 7 chars, which the old `len > 7 && addr[0:6] == "tcp://"`
// check skipped, leaving the scheme in place.
func TestGrpcConfigShortPrefix(t *testing.T) {
	cfg := &Config{GrpcAPIAddress: "tcp://x"}
	require.NoError(t, cfg.InitDefaults())

	assert.Equal(t, "x", cfg.GrpcAPIAddress)
}

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{}
	require.NoError(t, cfg.InitDefaults())

	assert.Equal(t, "127.0.0.1:10000", cfg.GrpcAPIAddress)
	assert.Equal(t, "tcp://127.0.0.1:30000", cfg.ProxyAddress)
	assert.Equal(t, "roadrunner", cfg.Name)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.NotNil(t, cfg.Pool)
}

func TestConfigTLSMissingKey(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{TLS: &TLS{
		Key:  filepath.Join(dir, "absent.key"),
		Cert: filepath.Join(dir, "absent.cert"),
	}}

	err := cfg.InitDefaults()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key file")
	assert.Contains(t, err.Error(), "does not exists")
}

func TestConfigTLSMissingCert(t *testing.T) {
	dir := t.TempDir()
	key := filepath.Join(dir, "tls.key")
	require.NoError(t, os.WriteFile(key, []byte("x"), 0o600))

	cfg := &Config{TLS: &TLS{Key: key, Cert: filepath.Join(dir, "absent.cert")}}

	err := cfg.InitDefaults()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cert file")
	assert.Contains(t, err.Error(), "does not exists")
}

func TestConfigTLSValid(t *testing.T) {
	dir := t.TempDir()
	key := filepath.Join(dir, "tls.key")
	cert := filepath.Join(dir, "tls.cert")
	require.NoError(t, os.WriteFile(key, []byte("x"), 0o600))
	require.NoError(t, os.WriteFile(cert, []byte("x"), 0o600))

	cfg := &Config{TLS: &TLS{Key: key, Cert: cert}}

	require.NoError(t, cfg.InitDefaults())
}
