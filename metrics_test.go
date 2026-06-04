package centrifuge

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/roadrunner-server/pool/v2/fsm"
	"github.com/roadrunner-server/pool/v2/state/process"
	"github.com/stretchr/testify/require"
)

type fakeInformer struct {
	states []*process.State
}

func (f *fakeInformer) Workers() []*process.State { return f.states }

// collectCount registers exp in a fresh registry and returns the number of
// gathered samples. Gather rejects duplicate label sets, which is why the
// states below must carry distinct PIDs.
func collectCount(t *testing.T, exp *StatsExporter) int {
	t.Helper()

	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(exp))

	mfs, err := reg.Gather()
	require.NoError(t, err)

	var n int
	for _, mf := range mfs {
		n += len(mf.GetMetric())
	}

	return n
}

func TestStatsExporterCollect(t *testing.T) {
	// Distinct PIDs avoid duplicate label sets; the three statuses also cover the
	// ready/working/default arms of the status switch in Collect.
	inf := &fakeInformer{states: []*process.State{
		{Pid: 1, Status: fsm.StateReady, StatusStr: "ready", MemoryUsage: 100},
		{Pid: 2, Status: fsm.StateWorking, StatusStr: "working", MemoryUsage: 200},
		{Pid: 3, Status: fsm.StateInactive, StatusStr: "inactive", MemoryUsage: 300},
	}}

	// 2 per-worker series (state + memory) plus 5 aggregate series.
	require.Equal(t, 2*len(inf.states)+5, collectCount(t, newWorkersExporter(inf)))
}

func TestStatsExporterCollectEmpty(t *testing.T) {
	require.Equal(t, 5, collectCount(t, newWorkersExporter(&fakeInformer{})))
}

func TestStatsExporterDescribe(t *testing.T) {
	exp := newWorkersExporter(&fakeInformer{})

	ch := make(chan *prometheus.Desc, 16)
	exp.Describe(ch)
	close(ch)

	var n int
	for range ch {
		n++
	}

	require.Equal(t, 7, n)
}

func TestPluginMetricsCollector(t *testing.T) {
	p := &Plugin{statsExporter: newWorkersExporter(&fakeInformer{})}

	require.Len(t, p.MetricsCollector(), 1)
}
