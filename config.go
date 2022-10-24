package centrifuge

import (
	"github.com/roadrunner-server/sdk/v3/pool"
)

type Config struct {
	ProxyAddress   string `mapstructure:"proxy_address"`
	Endpoint       string `mapstructure:"endpoint"`
	Name           string `mapstructure:"name"`
	Version        string `mapstructure:"version"`
	TokenSecretKey string `mapstructure:"token_hmac_secret_key"`

	Pool *pool.Config `mapstructure:"pool"`
}

func (c *Config) InitDefaults() {
	if c.Endpoint == "" {
		c.Endpoint = "ws://localhost:8000/connection/websocket"
	}

	if c.ProxyAddress == "" {
		c.ProxyAddress = "tcp://127.0.0.1:30000"
	}

	if c.Name == "" {
		c.Name = "roadrunner"
	}

	if c.Version == "" {
		c.Version = "1.0.0"
	}

	if c.TokenSecretKey == "" {
		c.TokenSecretKey = "test"
	}

	if c.Pool == nil {
		c.Pool = &pool.Config{}
		c.Pool.InitDefaults()
	}
}
