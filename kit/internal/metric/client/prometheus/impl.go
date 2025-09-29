package colletor

import (
	"context"
	"strings"
	"sync"

	"github.com/KingTrack/gin-kit/kit/types/metric/conf"

	"github.com/KingTrack/gin-kit/kit/internal/metric/client/label"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

type Collector struct {
	metrics   metrics.Registry
	container *label.Container
	descCache map[string]*prometheus.Desc
	mu        sync.RWMutex
}

func New(metrics metrics.Registry, container *label.Container) *Collector {
	return &Collector{
		metrics:   metrics,
		container: container,
		descCache: make(map[string]*prometheus.Desc),
	}
}

func (c *Collector) Init(ctx context.Context, config *conf.Config) error {
	if err := prometheus.NewRegistry().Register(c); err != nil {
		return err
	}
	return nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	// 发送一个通用的无效描述符
	// 这样Prometheus会跳过预检查，直接调用Collect
	ch <- prometheus.NewInvalidDesc(nil)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.metrics.Each(func(name string, metric interface{}) {
		metricInfo := c.container.GetMetricInfo(name)

		labelKeys := metricInfo.GetLabelKeys()
		labelValues := metricInfo.GetLabelValues()

		switch m := metric.(type) {
		case metrics.Counter:
			desc := c.getOrCreateDesc(metricInfo.GetCounterTotalMetricName(), "Counter metric", labelKeys)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(m.Count()), labelValues...)
		case metrics.Gauge:
			desc := c.getOrCreateDesc(metricInfo.GetGaugeMetricName(), "Gauge metric", labelKeys)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(m.Value()), labelValues...)
		case metrics.GaugeFloat64:
			desc := c.getOrCreateDesc(metricInfo.GetGaugeMetricName(), "Float64 gauge metric", labelKeys)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, m.Value(), labelValues...)
		case metrics.Meter:
			countDesc := c.getOrCreateDesc(metricInfo.GetMeterCountMetricName(), "Meter event count", labelKeys)
			rate1Desc := c.getOrCreateDesc(metricInfo.GetMeterRate1MetricName(), "Meter 1-minute rate (QPS)", labelKeys)
			rateMeanDesc := c.getOrCreateDesc(metricInfo.GetMeterRateMeanMetricName(), "Meter mean rate (QPS)", labelKeys)

			ch <- prometheus.MustNewConstMetric(countDesc, prometheus.GaugeValue, float64(m.Count()), labelValues...)
			ch <- prometheus.MustNewConstMetric(rate1Desc, prometheus.GaugeValue, m.Rate1(), labelValues...)
			ch <- prometheus.MustNewConstMetric(rateMeanDesc, prometheus.GaugeValue, m.RateMean(), labelValues...)
		case metrics.Timer:
			// 统计指标：Count, Mean, Min, Max（基础性能指标）
			countDesc := c.getOrCreateDesc(metricInfo.GetTimerCountMetricName(), "Timer event count", labelKeys)                    // 请求总数 (Count)
			meanDesc := c.getOrCreateDesc(metricInfo.GetTimerMeanSecondsMetricName(), "Timer mean duration in seconds", labelKeys)  // 平均值：平均耗时（纳秒）
			minDesc := c.getOrCreateDesc(metricInfo.GetTimerMaxSecondsMetricName(), "Timer minimum duration in seconds", labelKeys) // 最大值：最长耗时（纳秒）
			maxDesc := c.getOrCreateDesc(metricInfo.GetTimerMinSecondsMetricName(), "Timer maximum duration in seconds", labelKeys) // 最小值：最短耗时（纳秒）

			// 速率指标：Rate1, RateMean（QPS监控）
			rate1Desc := c.getOrCreateDesc(metricInfo.GetTimerRate1MetricName(), "Timer 1-minute rate (QPS)", labelKeys)   // 1分钟移动平均 QPS（每秒请求数）
			rateMeanDesc := c.getOrCreateDesc(metricInfo.GetTimerRateMeanMetricName(), "Timer mean rate (QPS)", labelKeys) // 平均 QPS（每秒请求数）

			// 分位数指标：P50, P95, P99（SLA监控）
			p50Desc := c.getOrCreateDesc(metricInfo.GetTimerP50SecondsMetricName(), "Timer 50th percentile in seconds", labelKeys)    // P50响应时间 (Percentile(0.5))
			p95Desc := c.getOrCreateDesc(metricInfo.GetTimerP95SecondsMetricName(), "Timer 95th percentile in seconds", labelKeys)    // P95响应时间 (Percentile(0.95))
			p99Desc := c.getOrCreateDesc(metricInfo.GetTimerP99SecondsMetricName(), "Timer 99th percentile in seconds", labelKeys)    // P99响应时间 (Percentile(0.99))
			p999Desc := c.getOrCreateDesc(metricInfo.GetTimerP999SecondsMetricName(), "Timer 999th percentile in seconds", labelKeys) // P99响应时间 (Percentile(0.999))

			ch <- prometheus.MustNewConstMetric(countDesc, prometheus.GaugeValue, float64(m.Count()), labelValues...)
			ch <- prometheus.MustNewConstMetric(meanDesc, prometheus.GaugeValue, m.Mean()/1e9, labelValues...)
			ch <- prometheus.MustNewConstMetric(minDesc, prometheus.GaugeValue, float64(m.Min())/1e9, labelValues...)
			ch <- prometheus.MustNewConstMetric(maxDesc, prometheus.GaugeValue, float64(m.Max())/1e9, labelValues...)

			ch <- prometheus.MustNewConstMetric(rate1Desc, prometheus.GaugeValue, m.Rate1(), labelValues...)
			ch <- prometheus.MustNewConstMetric(rateMeanDesc, prometheus.GaugeValue, m.RateMean(), labelValues...)

			ch <- prometheus.MustNewConstMetric(p50Desc, prometheus.GaugeValue, m.Percentile(0.5)/1e9, labelValues...)
			ch <- prometheus.MustNewConstMetric(p95Desc, prometheus.GaugeValue, m.Percentile(0.95)/1e9, labelValues...)
			ch <- prometheus.MustNewConstMetric(p99Desc, prometheus.GaugeValue, m.Percentile(0.99)/1e9, labelValues...)
			ch <- prometheus.MustNewConstMetric(p999Desc, prometheus.GaugeValue, m.Percentile(0.999)/1e9, labelValues...)
		case metrics.Histogram:
			// 核心指标
			countDesc := c.getOrCreateDesc(metricInfo.GetHistogramCountMetricName(), "Histogram sample count", labelKeys)
			meanDesc := c.getOrCreateDesc(metricInfo.GetHistogramMeanMetricName(), "Histogram mean value", labelKeys)
			minDesc := c.getOrCreateDesc(metricInfo.GetHistogramMinMetricName(), "Histogram minimum value", labelKeys)
			maxDesc := c.getOrCreateDesc(metricInfo.GetHistogramMaxMetricName(), "Histogram maximum value", labelKeys)

			// 核心分位数指标
			p50Desc := c.getOrCreateDesc(metricInfo.GetHistogramP50MetricName(), "Histogram 50th percentile", labelKeys)
			p95Desc := c.getOrCreateDesc(metricInfo.GetHistogramP95MetricName(), "Histogram 95th percentile", labelKeys)
			p99Desc := c.getOrCreateDesc(metricInfo.GetHistogramP99MetricName(), "Histogram 99th percentile", labelKeys)
			p999Desc := c.getOrCreateDesc(metricInfo.GetHistogramP999MetricName(), "Histogram 999th percentile", labelKeys)

			// 上报 Histogram 指标
			ch <- prometheus.MustNewConstMetric(countDesc, prometheus.GaugeValue, float64(m.Count()), labelValues...)
			ch <- prometheus.MustNewConstMetric(meanDesc, prometheus.GaugeValue, m.Mean(), labelValues...)
			ch <- prometheus.MustNewConstMetric(minDesc, prometheus.GaugeValue, float64(m.Min()), labelValues...)
			ch <- prometheus.MustNewConstMetric(maxDesc, prometheus.GaugeValue, float64(m.Max()), labelValues...)

			ch <- prometheus.MustNewConstMetric(p50Desc, prometheus.GaugeValue, m.Percentile(0.5), labelValues...)
			ch <- prometheus.MustNewConstMetric(p95Desc, prometheus.GaugeValue, m.Percentile(0.95), labelValues...)
			ch <- prometheus.MustNewConstMetric(p99Desc, prometheus.GaugeValue, m.Percentile(0.99), labelValues...)
			ch <- prometheus.MustNewConstMetric(p999Desc, prometheus.GaugeValue, m.Percentile(0.999), labelValues...)
		}
	})
}

func (c *Collector) getOrCreateDesc(name, help string, labelKeys []string) *prometheus.Desc {
	cacheKey := getDescCacheKey(name, labelKeys)

	if desc := c.getDesc(cacheKey); desc != nil {
		return desc
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if desc := c.descCache[cacheKey]; desc != nil {
		return desc
	}

	desc := prometheus.NewDesc(name, help, labelKeys, nil)
	c.descCache[cacheKey] = desc
	return desc
}

func (c *Collector) getDesc(cacheKey string) *prometheus.Desc {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.descCache[cacheKey]
}

func getDescCacheKey(name string, labelKeys []string) string {
	builder := &strings.Builder{}
	builder.WriteString(name)
	for _, v := range labelKeys {
		builder.WriteString(v)
	}
	return builder.String()
}
