package metrics

import (
	"expvar"
	"fmt"
	"os"
	"time"

	logging "github.com/op/go-logging"
	"github.com/rcrowley/go-metrics"
)

// InternalMetrics holds all metrics counters and timers.
type InternalMetrics struct {
	PID     *expvar.Int
	Version *expvar.String
	UpSince *expvar.String
	DBPath  *expvar.String

	CurrItems        metrics.Counter
	TotalItems       metrics.Counter
	TotalConnections metrics.Counter
	TotalThreads     metrics.Counter
	CurrThreads      metrics.Counter
	CmdGet           metrics.Counter
	CmdSet           metrics.Counter
	GetHits          metrics.Counter
	GetMisses        metrics.Counter
	ProtocolErrors   metrics.Counter
	NetworkErrors    metrics.Counter
	ReadonlyErrors   metrics.Counter
	ResponseTiming   metrics.Timer
}

// NewInternalMetrics initializes and registers all metrics.
func NewInternalMetrics(dbp string, dumpLogs bool) *InternalMetrics {
	m := &InternalMetrics{
		PID:              expvar.NewInt("pid"),
		Version:          expvar.NewString("version"),
		UpSince:          expvar.NewString("up_since"),
		DBPath:           expvar.NewString("db_path"),
		CurrItems:        metrics.NewCounter(),
		TotalItems:       metrics.NewCounter(),
		TotalConnections: metrics.NewCounter(),
		TotalThreads:     metrics.NewCounter(),
		CurrThreads:      metrics.NewCounter(),
		CmdGet:           metrics.NewCounter(),
		CmdSet:           metrics.NewCounter(),
		GetHits:          metrics.NewCounter(),
		GetMisses:        metrics.NewCounter(),
		ProtocolErrors:   metrics.NewCounter(),
		NetworkErrors:    metrics.NewCounter(),
		ReadonlyErrors:   metrics.NewCounter(),
		ResponseTiming:   metrics.NewTimer(),
	}

	m.PID.Set(int64(os.Getpid()))
	m.Version.Set("BEANO Server")
	m.UpSince.Set(time.Now().Format(time.RFC3339))
	m.DBPath.Set(dbp)

	metrics.Register("current_items", m.CurrItems)
	metrics.Register("total_items", m.TotalItems)
	metrics.Register("total_connections", m.TotalConnections)
	metrics.Register("total_threads", m.TotalThreads)
	metrics.Register("curr_threads", m.CurrThreads)
	metrics.Register("cmd_get", m.CmdGet)
	metrics.Register("cmd_set", m.CmdSet)
	metrics.Register("get_hits", m.GetHits)
	metrics.Register("get_misses", m.GetMisses)
	metrics.Register("protocol_errors", m.ProtocolErrors)
	metrics.Register("network_errors", m.NetworkErrors)
	metrics.Register("readonly_errors", m.ReadonlyErrors)
	metrics.Register("response_timing", m.ResponseTiming)

	if dumpLogs {
		go metrics.Log(metrics.DefaultRegistry, 60*time.Second, logging.NewLogBackend(os.Stdout, "", 0).Logger)
	}

	return m
}

// Metrics2Expvar publishes metrics from the registry to expvar.
func Metrics2Expvar(r metrics.Registry) {
	du := float64(time.Nanosecond)
	percentiles := []float64{0.50, 0.75, 0.95, 0.99, 0.999}
	r.Each(func(name string, i interface{}) {
		switch m := i.(type) {
		case metrics.Counter:
			expvar.Publish(fmt.Sprintf("%s.counter", name), expvar.Func(func() interface{} { return m.Count() }))
		case metrics.Meter:
			expvar.Publish(fmt.Sprintf("%s.rate.1", name), expvar.Func(func() interface{} { return m.Rate1() }))
			expvar.Publish(fmt.Sprintf("%s.rate.5", name), expvar.Func(func() interface{} { return m.Rate5() }))
			expvar.Publish(fmt.Sprintf("%s.rate.15", name), expvar.Func(func() interface{} { return m.Rate15() }))
			expvar.Publish(fmt.Sprintf("%s.rate.mean", name), expvar.Func(func() interface{} { return m.RateMean() }))
		case metrics.Histogram:
			expvar.Publish(fmt.Sprintf("%s.count", name), expvar.Func(func() interface{} { return m.Count() }))
			expvar.Publish(fmt.Sprintf("%s.mean", name), expvar.Func(func() interface{} { return m.Mean() }))
			expvar.Publish(fmt.Sprintf("%s.min", name), expvar.Func(func() interface{} { return m.Min() }))
			expvar.Publish(fmt.Sprintf("%s.max", name), expvar.Func(func() interface{} { return m.Max() }))
			expvar.Publish(fmt.Sprintf("%s.stddev", name), expvar.Func(func() interface{} { return m.StdDev() }))
			expvar.Publish(fmt.Sprintf("%s.variance", name), expvar.Func(func() interface{} { return m.Variance() }))
			for _, p := range percentiles {
				pct := p
				expvar.Publish(fmt.Sprintf("%s.percentile.%2.3f", name, pct), expvar.Func(func() interface{} { return m.Percentile(pct) }))
			}
		case metrics.Timer:
			expvar.Publish(fmt.Sprintf("%s.rate.1", name), expvar.Func(func() interface{} { return m.Rate1() }))
			expvar.Publish(fmt.Sprintf("%s.rate.5", name), expvar.Func(func() interface{} { return m.Rate5() }))
			expvar.Publish(fmt.Sprintf("%s.rate.15", name), expvar.Func(func() interface{} { return m.Rate15() }))
			expvar.Publish(fmt.Sprintf("%s.rate.mean", name), expvar.Func(func() interface{} { return m.RateMean() }))
			expvar.Publish(fmt.Sprintf("%s.mean", name), expvar.Func(func() interface{} { return du * m.Mean() }))
			expvar.Publish(fmt.Sprintf("%s.min", name), expvar.Func(func() interface{} { return int64(du) * m.Min() }))
			expvar.Publish(fmt.Sprintf("%s.max", name), expvar.Func(func() interface{} { return int64(du) * m.Max() }))
			expvar.Publish(fmt.Sprintf("%s.stddev", name), expvar.Func(func() interface{} { return du * m.StdDev() }))
			expvar.Publish(fmt.Sprintf("%s.variance", name), expvar.Func(func() interface{} { return du * m.Variance() }))
			for _, p := range percentiles {
				pct := p
				expvar.Publish(fmt.Sprintf("%s.percentile.%2.3f", name, pct), expvar.Func(func() interface{} { return m.Percentile(pct) }))
			}
		}
	})
	expvar.Publish("time", expvar.Func(func() interface{} { return time.Now().Format(time.RFC3339Nano) }))
}
