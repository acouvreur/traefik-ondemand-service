package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Status string

const (
	Starting Status = "Starting"
	Started         = "Started"
	Unknown         = "Unknown"
)

var (
	serviceStatusMetric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "service",
			Name:      "status",
			Help:      "Indicates if a service is starting, started or in an unknown state. When the service is scaled down this metric is deleted",
		},
		[]string{"service_name", "status"},
	)
	serviceScaleDownErrorMetric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "service",
			Name:      "scale_down_error",
			Help:      "Indicates whether a service scale down has failed",
		},
		[]string{"service_name"},
	)
	scaleDownErrorMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "scaler",
			Name:      "scale_down_errors",
			Help:      "Indicates how many scale down errors occured",
		},
	)
	serviceLastStartedTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "service",
			Name:      "last_started_time",
			Help:      "Indicates when the service was last started",
		},
		[]string{"service_name"},
	)
	serviceTimeout = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "service",
			Name:      "timeout",
			Help:      "Indicates the duration after which the service will be scaled down",
		},
		[]string{"service_name"},
	)
)

func getCurriedMetricsFromGaugeVecByServiceName(metricToCurry *prometheus.GaugeVec, serviceName string) *prometheus.GaugeVec {
	return serviceStatusMetric.MustCurryWith(prometheus.Labels{
		"service_name": serviceName,
	})
}
func getMetricFromGaugeVecByServiceName(metric *prometheus.GaugeVec, serviceName string) prometheus.Gauge {
	return metric.With(prometheus.Labels{
		"service_name": serviceName,
	})
}

func setStatusInMetric(metric *prometheus.GaugeVec, status Status) {
	for _, s := range []Status{Starting, Started, Unknown} {
		if status == s {
			metric.With(prometheus.Labels{"status": string(s)}).Set(1)
		} else {
			metric.With(prometheus.Labels{"status": string(s)}).Set(0)
		}
	}
}

func OnStoreUpdate(serviceName string, storeTimeout time.Duration) {
	getMetricFromGaugeVecByServiceName(serviceLastStartedTime, serviceName).SetToCurrentTime()
	getMetricFromGaugeVecByServiceName(serviceTimeout, serviceName).Set(float64(storeTimeout))
}

func OnServiceStateChange(serviceName string, status Status) {
	setStatusInMetric(getCurriedMetricsFromGaugeVecByServiceName(serviceStatusMetric, serviceName), status)
}

func OnScaleDown(serviceName string) {
	serviceStatusMetric.Delete(prometheus.Labels{
		"service_name": serviceName,
	})
}

func OnScaleDownError(serviceName string) {
	scaleDownErrorMetric.Inc()
	getMetricFromGaugeVecByServiceName(serviceScaleDownErrorMetric, serviceName).Set(1)
}
