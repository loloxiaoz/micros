package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

func initPrometheus() {
	HttpUrlStat = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "micros_http_request_total",
			Help: "How many HTTP requests processed, partitioned by status code and url",
		},
		[]string{"code", "url"},
	)
	prometheus.MustRegister(HttpUrlStat)

	HttpTimeStat = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "micros_http_request_latency",
			Help:       "time latency of HTTP requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"label"},
	)
	prometheus.MustRegister(HttpTimeStat)
}
