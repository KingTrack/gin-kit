package label

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type MetricInfo struct {
	Name   string  `json:"name"`
	Labels []Label `json:"labels"`
}

func (m *MetricInfo) GetLabelKeys() []string {
	var keys []string
	for _, v := range m.Labels {
		keys = append(keys, v.Key)
	}
	return keys
}

func (m *MetricInfo) GetLabelValues() []string {
	var values []string
	for _, v := range m.Labels {
		values = append(values, v.Value)
	}
	return values
}

type Label struct {
	Key   string
	Value string
}

func (l Label) ToString() string {
	return strings.Join([]string{l.Key, l.Value}, "=")
}

type Container struct {
	cache sync.Map
}

func New() *Container {
	return &Container{}
}

func (c *Container) GetMetricInfo(fullName string) *MetricInfo {
	val, ok := c.cache.Load(fullName)
	if !ok {
		return nil
	}
	x, ok := val.(*MetricInfo)
	if !ok {
		return nil
	}
	return x
}

func (c *Container) GetMetricName(baseName string, labels []Label) string {
	builder := strings.Builder{}
	builder.WriteString(baseName)
	for _, v := range labels {
		builder.WriteString(v.ToString())
	}
	return builder.String()
}

func (c *Container) RegisterMetricInfo(fullName string, baseName string, labels []Label) {
	_, ok := c.cache.Load(fullName)
	if ok {
		return
	}

	c.cache.Store(fullName, &MetricInfo{
		Name:   baseName,
		Labels: labels,
	})
}

func ParseLabels(kvs ...interface{}) []Label {
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		return nil
	}

	var labels []Label
	// key 只支持字符串, value 只支持 int bool string 3种基础类型，其他类型不支持
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			return nil
		}

		var valueStr string
		value := kvs[i+1]
		switch v := value.(type) {
		case string:
			valueStr = v
		case bool:
			valueStr = fmt.Sprintf("%t", v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			valueStr = fmt.Sprintf("%d", v)
		default:
			return nil
		}

		labels = append(labels, Label{Key: key, Value: valueStr})
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Key < labels[j].Key
	})

	return labels
}
