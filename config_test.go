package centrifuge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrpcConfig(t *testing.T) {
	cfg := &Config{GrpcApiAddress: "tcp://foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "foo.bar", cfg.GrpcApiAddress)
}

func TestGrpcConfig1(t *testing.T) {
	cfg := &Config{GrpcApiAddress: "foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "foo.bar", cfg.GrpcApiAddress)
}

func TestGrpcConfig2(t *testing.T) {
	cfg := &Config{GrpcApiAddress: "tcp:/foo.bar"}
	err := cfg.InitDefaults()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "tcp:/foo.bar", cfg.GrpcApiAddress)
}
