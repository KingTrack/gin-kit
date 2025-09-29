package pusher

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/KingTrack/gin-kit/kit/types/metric/conf"

	"github.com/KingTrack/gin-kit/kit/internal/metric/client/label"
	"github.com/KingTrack/gin-kit/kit/internal/metric/client/n9e/client"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/rcrowley/go-metrics"
)

type Pusher struct {
	config    *conf.Config
	registry  metrics.Registry
	container *label.Container
	client    *client.Client
}

func New(registry metrics.Registry, container *label.Container) *Pusher {
	return &Pusher{
		registry:  registry,
		container: container,
	}
}

func (p *Pusher) Init(ctx context.Context, config *conf.Config) error {
	p.config = config

	p.client = client.New(
		config.N9e.URL,
		client.WithToken(config.N9e.Token),
		client.WithTimeout(time.Second*10),
	)

	go p.run()
	return nil
}

func (p *Pusher) run() {
	interval := time.Duration(p.config.N9e.IntervalSec) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.send(context.Background()); err != nil {
				runtime.Get().LoggerRegistry().GenLogger().Printf("n9e pusher send metrics error:%v", err)
			}
		}
	}
}

func (p *Pusher) send(ctx context.Context) error {
	now := time.Now().Unix()
	step := p.config.N9e.StepSec
	var n9eMetrics []client.Metric

	p.registry.Each(func(name string, metric interface{}) {
		metricInfo := p.container.GetMetricInfo(name)
		if metricInfo == nil {
			return
		}

		tags := buildTags(metricInfo.Labels)

		switch m := metric.(type) {
		case metrics.Counter:
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetCounterTotalMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(m.Count()),
				CounterType: client.CounterTypeCOUNTER,
				Tags:        tags,
			})

		case metrics.Gauge:
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetGaugeMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(m.Value()),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})

		case metrics.GaugeFloat64:
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetGaugeMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       m.Value(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})

		case metrics.Meter:
			snapshot := m.Snapshot()
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetMeterCountMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Count()),
				CounterType: client.CounterTypeCOUNTER,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetMeterRate1MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Rate1(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetMeterRateMeanMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.RateMean(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})

		case metrics.Timer:
			snapshot := m.Snapshot()
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerCountMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Count()),
				CounterType: client.CounterTypeCOUNTER,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerMeanSecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Mean() / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerMinSecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Min()) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerMaxSecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Max()) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerRate1MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Rate1(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerRateMeanMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.RateMean(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerP50SecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.5) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerP95SecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.95) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerP99SecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.99) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetTimerP999SecondsMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.999) / 1e9,
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})

		case metrics.Histogram:
			snapshot := m.Snapshot()
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramCountMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Count()),
				CounterType: client.CounterTypeCOUNTER,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramMeanMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Mean(),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramMinMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Min()),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramMaxMetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       float64(snapshot.Max()),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramP50MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.5),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramP95MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.95),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramP99MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.99),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
			n9eMetrics = append(n9eMetrics, client.Metric{
				Endpoint:    p.config.Endpoint,
				Metric:      metricInfo.GetHistogramP999MetricName(),
				Timestamp:   now,
				Step:        step,
				Value:       snapshot.Percentile(0.999),
				CounterType: client.CounterTypeGAUGE,
				Tags:        tags,
			})
		}
	})

	return p.client.PushMetrics(n9eMetrics)
}

func buildTags(labels []label.Label) string {
	if len(labels) == 0 {
		return ""
	}

	pairs := make([]string, 0, len(labels))
	for _, v := range labels {
		pairs = append(pairs, v.ToString())
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
}
