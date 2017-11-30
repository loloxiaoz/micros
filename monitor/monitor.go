package monitor

import (
	"github.com/getsentry/raven-go"
	"github.com/prometheus/client_golang/prometheus"
)

var client *raven.Client
var HttpUrlStat *prometheus.CounterVec
var HttpTimeStat *prometheus.SummaryVec

func init() {
	initSentry()
	initPrometheus()
}
