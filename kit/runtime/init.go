package runtime

import (
	"sync"

	"github.com/KingTrack/gin-kit/kit/engine"
)

var (
	runtimeE *engine.Engine
	onceE    sync.Once
)

func Set(e *engine.Engine) {
	if e == nil {
		panic("engine cannot be nil")
	}
	onceE.Do(func() {
		runtimeE = e
	})
}

func Get() *engine.Engine {
	if runtimeE == nil {
		onceE.Do(func() {
			runtimeE = engine.NewDefault()
		})
	}
	return runtimeE
}
