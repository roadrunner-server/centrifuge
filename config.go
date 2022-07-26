package centrifuge

import (
	"os"

	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/sdk/v3/pool"
)

type Config struct {
	// host + port
	ProxyAddress string `mapstructure:"proxy_address"`
	// host + port
	GrpcApiAddress string `mapstructure:"grpc_api_address"`
	UseCompressor  bool   `mapstructure:"use_compressor"`
	Version        string `mapstructure:"version"`
	Name           string `mapstructure:"name"`
	TLS            *TLS   `mapstructure:"tls"`

	Pool *pool.Config `mapstructure:"pool"`
}

type TLS struct {
	Key  string `mapstructure:"key"`
	Cert string `mapstructure:"cert"`
}

func (c *Config) InitDefaults() error {
	const op = errors.Op("centrifuge_init_defaults")

	if c.GrpcApiAddress == "" {
		c.GrpcApiAddress = "127.0.0.1:10000"
	}

	if len(c.GrpcApiAddress) > 7 && c.GrpcApiAddress[0:6] == "tcp://" {
		c.GrpcApiAddress = c.GrpcApiAddress[6:]
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

	if c.Pool == nil {
		c.Pool = &pool.Config{}
		c.Pool.InitDefaults()
	}

	if c.TLS != nil {
		if _, err := os.Stat(c.TLS.Key); err != nil {
			if os.IsNotExist(err) {
				return errors.E(op, errors.Errorf("key file '%s' does not exists", c.TLS.Key))
			}

			return errors.E(op, err)
		}

		if _, err := os.Stat(c.TLS.Cert); err != nil {
			if os.IsNotExist(err) {
				return errors.E(op, errors.Errorf("cert file '%s' does not exists", c.TLS.Cert))
			}

			return errors.E(op, err)
		}
	}

	return nil
}
