package label

// Counter相关指标名称
func (m *MetricInfo) GetCounterTotalMetricName() string {
	return m.Name + "_total"
}

// Gauge相关指标名称
func (m *MetricInfo) GetGaugeMetricName() string {
	return m.Name
}

// Meter相关指标名称
func (m *MetricInfo) GetMeterCountMetricName() string {
	return m.Name + "_count"
}

func (m *MetricInfo) GetMeterRate1MetricName() string {
	return m.Name + "_rate1"
}

func (m *MetricInfo) GetMeterRateMeanMetricName() string {
	return m.Name + "_rate_mean"
}

// Timer相关指标名称
func (m *MetricInfo) GetTimerCountMetricName() string {
	return m.Name + "_count"
}

func (m *MetricInfo) GetTimerMeanSecondsMetricName() string {
	return m.Name + "_mean_seconds"
}

func (m *MetricInfo) GetTimerMinSecondsMetricName() string {
	return m.Name + "_min_seconds"
}

func (m *MetricInfo) GetTimerMaxSecondsMetricName() string {
	return m.Name + "_max_seconds"
}

func (m *MetricInfo) GetTimerRate1MetricName() string {
	return m.Name + "_rate1"
}

func (m *MetricInfo) GetTimerRateMeanMetricName() string {
	return m.Name + "_rate_mean"
}

func (m *MetricInfo) GetTimerP50SecondsMetricName() string {
	return m.Name + "_p50_seconds"
}

func (m *MetricInfo) GetTimerP95SecondsMetricName() string {
	return m.Name + "_p95_seconds"
}

func (m *MetricInfo) GetTimerP99SecondsMetricName() string {
	return m.Name + "_p99_seconds"
}

func (m *MetricInfo) GetTimerP999SecondsMetricName() string {
	return m.Name + "_p999_seconds"
}

// Histogram相关指标名称
func (m *MetricInfo) GetHistogramCountMetricName() string {
	return m.Name + "_count"
}

func (m *MetricInfo) GetHistogramMeanMetricName() string {
	return m.Name + "_mean"
}

func (m *MetricInfo) GetHistogramMinMetricName() string {
	return m.Name + "_min"
}

func (m *MetricInfo) GetHistogramMaxMetricName() string {
	return m.Name + "_max"
}

func (m *MetricInfo) GetHistogramP50MetricName() string {
	return m.Name + "_p50"
}

func (m *MetricInfo) GetHistogramP95MetricName() string {
	return m.Name + "_p95"
}

func (m *MetricInfo) GetHistogramP99MetricName() string {
	return m.Name + "_p99"
}

func (m *MetricInfo) GetHistogramP999MetricName() string {
	return m.Name + "_p999"
}
