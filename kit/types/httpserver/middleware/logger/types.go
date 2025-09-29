package logger

type IConfig interface {
	GetRequestStartTimeKey() string
	GetRequestDurationMsKey() string
	GetMethodNameValue(path string) string
	GetClientIPKey() string
	GetAppCodeKey() string
	GetHTTPStatusKey() string
}
