package conf

type BackendEnum string

const (
	BackendPrometheus BackendEnum = "prometheus"
	BackendN9e        BackendEnum = "n9e"
)

type Config struct {
	ServiceName string      `toml:"service_name" json:"service_name" yaml:"service_name"`
	Endpoint    string      `toml:"endpoint" json:"endpoint" yaml:"endpoint"`
	BackendName BackendEnum `toml:"backend_name" json:"backend_name" yaml:"backend_name"`
	N9e         N9e         `toml:"n9e" json:"n9e" yaml:"n9e"`
	Prometheus  Prometheus  `toml:"prometheus" json:"prometheus" yaml:"prometheus"`
}

type N9e struct {
	URL         string `toml:"url" json:"url" yaml:"url"`
	Token       string `toml:"token" json:"token" yaml:"token"`
	IntervalSec int64  `toml:"interval_sec" json:"interval_sec" yaml:"interval_sec"` // 上报间隔（秒）
	StepSec     int64  `toml:"step_sec" json:"step_sec" yaml:"step_sec"`             // 采集步长（秒）
}

type Prometheus struct {
	Path string `json:"path" json:"path" yaml:"path"`
}
