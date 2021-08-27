package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	createCounter prometheus.Counter
	updateCounter prometheus.Counter
	deleteCounter prometheus.Counter
)

func RegisterMetrics() {
	createCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "meeting_create_count_total",
		Help: "The total create meeting",
	})

	updateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "meeting_update_count_total",
		Help: "The total update meeting",
	})

	deleteCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "meeting_delete_count_total",
		Help: "The total delete meeting",
	})
}

func CreateCounterInc() {
	createCounter.Inc()
}

func UpdateCounterInc() {
	updateCounter.Inc()
}

func RemoveCounterInc() {
	deleteCounter.Inc()
}
