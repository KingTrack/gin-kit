package client

import (
	"time"

	"github.com/KingTrack/gin-kit/kit/internal/metric/client/label"
	"github.com/KingTrack/gin-kit/kit/runtime"
)

// 统计事件频率
func Meter(baseName string, value int64, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Meter(baseName, label.ParseLabels(kvs)).Mark(value)
}

// 值统计
func Gauge(baseName string, value int64, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Gauge(baseName, label.ParseLabels(kvs)).Update(value)
}

// 统计事件频率 + 耗时
func Timer(baseName string, since time.Time, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Timer(baseName, label.ParseLabels(kvs)).UpdateSince(since)
}

// 统计事件频率 + 耗时
func TimerDuration(baseName string, duration time.Duration, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Timer(baseName, label.ParseLabels(kvs)).Update(duration)
}

// 计数统计
func CounterInc(baseName string, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Counter(baseName, label.ParseLabels(kvs)).Inc(1)
}

// 计数统计
func CounterDec(baseName string, kvs ...interface{}) {
	runtime.Get().MetricRegistry().Counter(baseName, label.ParseLabels(kvs)).Dec(1)
}
