package registry

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/internal/metric/client/label"
	n9epusher "github.com/KingTrack/gin-kit/kit/internal/metric/client/n9e"
	prometheuscollector "github.com/KingTrack/gin-kit/kit/internal/metric/client/prometheus"
	"github.com/KingTrack/gin-kit/kit/types/metric/conf"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
)

type Registry struct {
	metrics   metrics.Registry
	container *label.Container
	config    *conf.Config
}

func New() *Registry {
	r := &Registry{
		metrics:   metrics.NewRegistry(),
		container: label.New(),
	}

	return r
}

func (r *Registry) Init(ctx context.Context, config *conf.Config) error {
	r.config = config

	switch config.BackendName {
	case conf.BackendPrometheus:
		prometheusCollector := prometheuscollector.New(r.metrics, r.container)
		if err := prometheusCollector.Init(ctx, config); err != nil {
			return err
		}
	case conf.BackendN9e:
		n9ePusher := n9epusher.New(r.metrics, r.container)
		if err := n9ePusher.Init(ctx, config); err != nil {
			return err
		}
	default:
		return errors.Errorf("unsupported backend: %s", config.BackendName)
	}

	return nil
}

// Gauge 创建或获取仪表盘
func (r *Registry) Gauge(baseName string, labels []label.Label) metrics.Gauge {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewGauge())

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.Gauge)
}

// GaugeFloat64 创建或获取浮点仪表盘
func (r *Registry) GaugeFloat64(baseName string, labels []label.Label) metrics.GaugeFloat64 {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewGaugeFloat64())

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.GaugeFloat64)
}

// Counter 创建或获取计数器
func (r *Registry) Counter(baseName string, labels []label.Label) metrics.Counter {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewCounter())

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.Counter)
}

// Meter 创建或获取计量器
func (r *Registry) Meter(baseName string, labels []label.Label) metrics.Meter {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewMeter())

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.Meter)
}

// Timer 创建或获取计时器
func (r *Registry) Timer(baseName string, labels []label.Label) metrics.Timer {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewTimer())

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.Timer)
}

// Histogram 创建或获取直方图
func (r *Registry) Histogram(baseName string, labels []label.Label) metrics.Histogram {
	fullName := r.container.GetMetricName(baseName, labels)

	metric := r.metrics.GetOrRegister(fullName, metrics.NewExpDecaySample(1028, 0.015))

	r.container.RegisterMetricInfo(fullName, baseName, labels)

	return metric.(metrics.Histogram)
}
