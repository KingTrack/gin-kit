package globals

import (
	contextregistry "github.com/KingTrack/gin-kit/kit/internal/context/registry"
	datacenterregistry "github.com/KingTrack/gin-kit/kit/internal/datacenter/registry"
	loggerregistry "github.com/KingTrack/gin-kit/kit/internal/logger/registry"
)

var (
	loggerRegistry     *loggerregistry.Registry
	contextRegistry    *contextregistry.Registry
	datacenterRegistry *datacenterregistry.Registry
)

func SetLogger(registry *loggerregistry.Registry) {
	loggerRegistry = registry
}

func GetLogger() *loggerregistry.Registry {
	return loggerRegistry
}

func SetContext(registry *contextregistry.Registry) {
	contextRegistry = registry
}

func GetContext() *contextregistry.Registry {
	return contextRegistry
}

func SetDatacenter(registry *datacenterregistry.Registry) {
	datacenterRegistry = registry
}

func GetDatacenter() *datacenterregistry.Registry {
	return datacenterRegistry
}
