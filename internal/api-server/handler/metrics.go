package handler

import (
	"net/http"

	"github.com/Donders-Institute/filer-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

// This trick of integrating promhttp handler with swagger server is taken from
// the blog: https://www.kaznacheev.me/posts/en/go_swagger_tricks/
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
func GetMetrics(ucache *UserResourceCache, pcache *ProjectResourceCache) func(p operations.GetMetricsParams) middleware.Responder {

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		userCount,
		projectCount,
		projectStorageQuota,
		projectStorageUsage,
	)

	log.Debugf("GetMetrics called %p", promRegistry)

	return func(p operations.GetMetricsParams) middleware.Responder {
		collectMetrics(ucache, pcache)
		return NewCustomResponder(
			p.HTTPRequest,
			promhttp.HandlerFor(
				promRegistry,
				promhttp.HandlerOpts{
					EnableOpenMetrics: false,
				},
			),
		)
	}
}

// metrics definition
var (
	projectCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "fgw_project_count",
			Help: "Total number of project storage directories",
		},
	)

	userCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "fgw_user_count",
			Help: "Total number of user home directories",
		},
	)

	projectStorageQuota = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fgw_project_storage_quota",
			Help: "The project storage quota in GiB",
		},
		[]string{
			// project number
			"number",
		},
	)

	projectStorageUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fgw_project_storage_usage",
			Help: "The project storage usage in GiB",
		},
		[]string{
			//project number
			"number",
		},
	)
)

// metrics recording function in an infinite loop
func collectMetrics(ucache *UserResourceCache, pcache *ProjectResourceCache) {
	// projects
	projectStorageQuota.Reset()
	projectStorageUsage.Reset()
	i := 0
	for pnumber, resc := range pcache.getAllResources(false) {
		i++
		projectStorageQuota.WithLabelValues(pnumber).Set(
			float64(*resc.storage.QuotaGb),
		)
		projectStorageUsage.WithLabelValues(pnumber).Set(
			float64(*resc.storage.UsageMb) / 1024,
		)
	}
	projectCount.Set(float64(i))

	// users
	userCount.Set(float64(len(ucache.getAllResources(false))))
}
