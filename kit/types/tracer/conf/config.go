package conf

type ProtoEnum string

const (
	ProtoOpenTelemetry ProtoEnum = "OpenTelemetry"
	ProtoOpenTracing   ProtoEnum = "OpenTracing"
)

type BackendEnum string

const (
	BackendJaeger     BackendEnum = "jaeger"
	BackendSkywalking BackendEnum = "skywalking" // 暂时不支持，太麻烦了，需要编译代码
	BackendZipkin     BackendEnum = "zipkin"
)

type Config struct {
	ServiceName string      `toml:"service_name" json:"service_name" yaml:"service_name"`
	Enabled     bool        `toml:"enabled" json:"enabled" yaml:"enabled"`
	ReportURL   string      `json:"report_url" json:"report_url" yaml:"report_url"`
	Proto       ProtoEnum   `toml:"proto" json:"proto" yaml:"proto"`
	BackendName BackendEnum `toml:"backend_name" json:"backend_name" yaml:"backend_name"`
}
