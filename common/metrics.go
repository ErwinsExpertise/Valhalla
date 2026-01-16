package common

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsPort metrics are served on
var MetricsPort string

var (
	// MetricsGauges available
	MetricsGauges = make(map[string]*prometheus.GaugeVec)

	// MetricsCounters available
	MetricsCounters = make(map[string]*prometheus.CounterVec)

	// metricsStarted ensures metrics server is only started once
	metricsStarted sync.Once
)

// StartMetrics initializes and handles metrics Prometheus endpoint
// This function is safe to call multiple times - it will only start the server once
func StartMetrics() {
	metricsStarted.Do(func() {
		go func() {
			http.Handle("/metrics", promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{},
			))
			if err := http.ListenAndServe("0.0.0.0:"+MetricsPort, nil); err != nil {
				log.Fatal(err)
			}
		}()
	})
}
