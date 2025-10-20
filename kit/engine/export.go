package engine

import (
	contextregistry "github.com/KingTrack/gin-kit/kit/internal/context/registry"
	datacenterregistry "github.com/KingTrack/gin-kit/kit/internal/datacenter/registry"
	httpclientregistry "github.com/KingTrack/gin-kit/kit/internal/httpclient/registry"
	loggerregistry "github.com/KingTrack/gin-kit/kit/internal/logger/registry"
	metricregistry "github.com/KingTrack/gin-kit/kit/internal/metric/registry"
	mysqlregistry "github.com/KingTrack/gin-kit/kit/internal/mysql/registry"
	redisregistry "github.com/KingTrack/gin-kit/kit/internal/redis/registry"
	tracerregistry "github.com/KingTrack/gin-kit/kit/internal/tracer/registry"
	serverconfig "github.com/KingTrack/gin-kit/kit/types/httpserver/conf"
	loggerconfig "github.com/KingTrack/gin-kit/kit/types/logger/conf"
)

func (e *Engine) MySQLRegistry() *mysqlregistry.Registry {
	return e.mysqlRegistry
}

func (e *Engine) RedisRegistry() *redisregistry.Registry {
	return e.redisRegistry
}

func (e *Engine) ServerConfig() *serverconfig.Config {
	return &e.config.Server
}

func (e *Engine) TracerRegistry() *tracerregistry.Registry {
	return e.tracerRegistry
}

func (e *Engine) LoggerRegistry() *loggerregistry.Registry {
	return e.loggerRegistry
}

func (e *Engine) LoggerConfig() *loggerconfig.Config {
	return &e.config.Logger
}

func (e *Engine) MetricRegistry() *metricregistry.Registry {
	return e.metricRegistry
}

func (e *Engine) ContextRegistry() *contextregistry.Registry {
	return e.contextRegistry
}

func (e *Engine) DatacenterRegistry() *datacenterregistry.Registry {
	return e.datacenterRegistry
}

func (e *Engine) HTTPClientRegistry() *httpclientregistry.Registry {
	return e.httpClientRegistry
}
