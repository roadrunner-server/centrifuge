module github.com/roadrunner-server/centrifuge/v6

go 1.26

toolchain go1.26.5

require (
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/prometheus/client_golang v1.24.0
	github.com/roadrunner-server/api-go/v6 v6.0.0-beta.13
	github.com/roadrunner-server/api-plugins/v6 v6.0.0-beta.2
	github.com/roadrunner-server/errors v1.5.0
	github.com/roadrunner-server/goridge/v4 v4.0.0-beta.3
	github.com/roadrunner-server/pool/v2 v2.0.0-beta.1
	github.com/roadrunner-server/tcplisten v1.5.2
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.82.1
	google.golang.org/protobuf v1.36.11
)

exclude (
	github.com/spf13/viper v1.18.0
	github.com/spf13/viper v1.18.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.0 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/roadrunner-server/events v1.0.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/tklauser/go-sysconf v0.4.0 // indirect
	github.com/tklauser/numcpus v0.12.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260720171339-e059f2f05d78 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
