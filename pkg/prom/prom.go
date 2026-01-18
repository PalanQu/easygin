package prom

import (
	"easygin/pkg/logging"
	"sync"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/zap"
)

const (
	DatabaseErrorTotal = "database_error_total"
)

type metricsType string

const (
	metricsTypeCounter   metricsType = "counter"
	metricsTypeGauge     metricsType = "gauge"
	metricsTypeHistogram metricsType = "histogram"
)

type customMetrics struct {
	MetricsName string
	Help        string
	Labels      []string
	MetricsType metricsType
}

var (
	metricsList = []customMetrics{
		{
			MetricsName: DatabaseErrorTotal,
			Help:        "Total number of database errors",
			Labels:      []string{"error_type"},
			MetricsType: metricsTypeCounter,
		},
	}

	instance *ginprom.Prometheus
	once     sync.Once
)

func Init(r *gin.Engine) *ginprom.Prometheus {
	once.Do(func() {
		registry := prometheus.NewRegistry()
		registry.MustRegister(collectors.NewGoCollector())
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		instance = ginprom.New(
			ginprom.Engine(r),
			ginprom.Subsystem("easygin"),
			ginprom.Path("/metrics"),
			ginprom.Registry(registry),
		)

		for _, metric := range metricsList {
			switch metric.MetricsType {
			case metricsTypeCounter:
				instance.AddCustomCounter(metric.MetricsName, metric.Help, metric.Labels)
			case metricsTypeGauge:
				instance.AddCustomGauge(metric.MetricsName, metric.Help, metric.Labels)
			case metricsTypeHistogram:
				instance.AddCustomHistogram(metric.MetricsName, metric.Help, metric.Labels)
			default:
				logging.GetGlobalLogger().Panic("unknown metrics type",
					zap.String("metrics_type", string(metric.MetricsType)))
			}
		}
	})

	return instance
}

func GetInstance() *ginprom.Prometheus {
	return instance
}
