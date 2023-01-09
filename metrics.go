package centrifuge

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/roadrunner-server/sdk/v4/metrics"
	"github.com/roadrunner-server/sdk/v4/state/process"
)

type Informer interface {
	Workers() []*process.State
}

func (p *Plugin) MetricsCollector() []prometheus.Collector {
	return []prometheus.Collector{p.statsExporter}
}

func newWorkersExporter(stats Informer) *metrics.StatsExporter {
	return &metrics.StatsExporter{
		TotalWorkersDesc: prometheus.NewDesc("rr_centrifugo_total_workers", "Total number of workers used by the Centrifugo plugin", nil, nil),
		TotalMemoryDesc:  prometheus.NewDesc("rr_centrifugo_workers_memory_bytes", "Memory usage by Centrifugo workers.", nil, nil),
		StateDesc:        prometheus.NewDesc("rr_centrifugo_worker_state", "Worker current state", []string{"state", "pid"}, nil),
		WorkerMemoryDesc: prometheus.NewDesc("rr_centrifugo_worker_memory_bytes", "Worker current memory usage", []string{"pid"}, nil),

		WorkersReady:   prometheus.NewDesc("rr_centrifugo_workers_ready", "Centrifugo workers currently in ready state", nil, nil),
		WorkersWorking: prometheus.NewDesc("rr_centrifugo_workers_working", "Centrifugo workers currently in working state", nil, nil),
		WorkersInvalid: prometheus.NewDesc("rr_centrifugo_workers_invalid", "Centrifugo workers currently in invalid,killing,destroyed,errored,inactive states", nil, nil),

		Workers: stats,
	}
}
