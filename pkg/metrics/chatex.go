package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var chatexCollector = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "autotrade",
	Subsystem: "chatex",
	Name:      "collector_total",
}, []string{"status"})

func ChatexCollectorOk() {
	chatexCollector.WithLabelValues("ok").Inc()
}

func ChatexCollectorErr() {
	chatexCollector.WithLabelValues("err").Inc()
}
