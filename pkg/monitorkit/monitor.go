package monitorkit

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "chatgw_gtp_requests_total",
		Help: "Total number of HTTP requests by status code and method.",
	},
	[]string{"skey"},
)
