package centrifuge

import (
	stderrors "errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/pool/v2/pool"
	"github.com/roadrunner-server/tcplisten"
)

type Config struct {
	// host + port
	ProxyAddress string                       `mapstructure:"proxy_address"`
	ProxySocket  *tcplisten.UnixSocketOptions `mapstructure:"proxy_socket"`
	// host + port
	GrpcAPIAddress string `mapstructure:"grpc_api_address"`
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

	if c.GrpcAPIAddress == "" {
		c.GrpcAPIAddress = "127.0.0.1:10000"
	}

	if addr, ok := strings.CutPrefix(c.GrpcAPIAddress, "tcp://"); ok {
		c.GrpcAPIAddress = addr
	}

	if c.ProxyAddress == "" {
		c.ProxyAddress = "tcp://127.0.0.1:30000"
	}

	if err := c.ProxySocket.Validate(c.ProxyAddress); err != nil {
		return errors.E(op, fmt.Errorf("centrifuge.proxy_socket: %w", err))
	}

	if c.Name == "" {
		c.Name = "roadrunner"
	}

	if c.Version == "" {
		c.Version = "1.0.0"
	}

	if c.Pool == nil {
		c.Pool = &pool.Config{}
	}
	c.Pool.InitDefaults()

	if c.TLS != nil { //nolint:nestif
		if _, err := os.Stat(c.TLS.Key); err != nil {
			if stderrors.Is(err, os.ErrNotExist) {
				return errors.E(op, errors.Errorf("key file '%s' does not exists", c.TLS.Key))
			}

			return errors.E(op, err)
		}

		if _, err := os.Stat(c.TLS.Cert); err != nil {
			if stderrors.Is(err, os.ErrNotExist) {
				return errors.E(op, errors.Errorf("cert file '%s' does not exists", c.TLS.Cert))
			}

			return errors.E(op, err)
		}
	}

	return nil
}

// Weak decoding into *int can convert booleans and fractions to IDs.
func validateUnixSocketIDs(cfg Configurer, key string) error {
	var raw map[string]any
	if err := cfg.UnmarshalKey(key, &raw); err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}

	const maxID = 4294967295
	for _, field := range []string{"uid", "gid"} {
		if raw[field] == nil {
			continue
		}
		value := reflect.ValueOf(raw[field])
		valid := false
		switch value.Kind() { //nolint:exhaustive // Other kinds are not valid IDs.
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			id := value.Int()
			valid = id >= 0 && id < maxID
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			valid = value.Uint() < maxID
		case reflect.String:
			id, err := strconv.ParseInt(value.String(), 0, strconv.IntSize)
			valid = err == nil && id >= 0 && id < maxID
		case reflect.Float32, reflect.Float64:
			id := value.Float()
			valid = id >= 0 && id < maxID && id == math.Trunc(id)
		}
		if !valid {
			return fmt.Errorf("%s.%s: must be an integer from 0 through 4294967294", key, field)
		}
	}
	return nil
}
