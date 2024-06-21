package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var HttpUrlStat *prometheus.CounterVec
var HttpTimeStat *prometheus.SummaryVec

func init() {
	initPrometheus()
}

func Report(flags map[string]string, err interface{}) {

}

