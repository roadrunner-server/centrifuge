package centrifuge

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
