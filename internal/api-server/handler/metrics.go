package handler

import (
	"net/http"
	"time"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}

func NewCustomResponder(r *http.Request, h http.Handler) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, _ runtime.Producer) {
		h.ServeHTTP(w, r)
	})
}

// GetMetrics handles the metrics request with the Prometheus handler
func GetMetrics() func(p operations.GetMetricsParams) middleware.Responder {

	recordMetrics()

	return func(p operations.GetMetricsParams) middleware.Responder {
		return NewCustomResponder(p.HTTPRequest, promhttp.Handler())
	}
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)
