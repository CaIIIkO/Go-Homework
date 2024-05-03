package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CustomizedCounterMetricAddPvz = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "add_pvz",
	})

	CustomizedCounterMetricGetById = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "GetById",
	})

	CustomizedCounterMetricUpdateById = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "UpdateById",
	})

	CustomizedCounterMetricDeleteById = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "DeleteById",
	})
)

func InitMetrics() {
	prometheus.MustRegister(CustomizedCounterMetricAddPvz)
	prometheus.MustRegister(CustomizedCounterMetricGetById)
	prometheus.MustRegister(CustomizedCounterMetricUpdateById)
	prometheus.MustRegister(CustomizedCounterMetricDeleteById)
}
