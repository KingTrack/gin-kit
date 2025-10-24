package logger

import (
	"sync"

	loggerregistry "github.com/KingTrack/gin-kit/kit/internal/logger/registry"
)

var (
	runtimeR *loggerregistry.Registry
	onceR    sync.Once
)

func Get() *loggerregistry.Registry {
	if runtimeR == nil {
		onceR.Do(func() {
			runtimeR = loggerregistry.New()
		})
	}
	return runtimeR
}

func Set(registry *loggerregistry.Registry) {
	runtimeR = registry
}
