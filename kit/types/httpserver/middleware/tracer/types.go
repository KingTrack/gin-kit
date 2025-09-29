package tracer

type IConfig interface {
	GetSpanName(method, uri string) string
	GetRequestDurationMsKey() string
	GetAppCodeKey() string
}
