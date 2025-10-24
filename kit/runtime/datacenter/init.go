package datacenter

import (
	"sync"

	datacenterregistry "github.com/KingTrack/gin-kit/kit/internal/datacenter/registry"
)

var (
	runtimeR *datacenterregistry.Registry
	onceR    sync.Once
)

func Set(registry *datacenterregistry.Registry) {
	runtimeR = registry
}

func Get() *datacenterregistry.Registry {
	if runtimeR == nil {
		onceR.Do(func() {
			runtimeR = datacenterregistry.New()
		})
	}
	return runtimeR
}
