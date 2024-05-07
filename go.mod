module github.com/roadrunner-server/centrifuge/v4

go 1.22.2

require (
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/goccy/go-json v0.10.2
	github.com/prometheus/client_golang v1.19.0
	github.com/roadrunner-server/api/v4 v4.12.0
	github.com/roadrunner-server/errors v1.4.0
	github.com/roadrunner-server/goridge/v3 v3.8.2
	github.com/roadrunner-server/sdk/v4 v4.7.2
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.1
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
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.14.0 // indirect
	github.com/roadrunner-server/tcplisten v1.4.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240506185236-b8a5c65736ae // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
