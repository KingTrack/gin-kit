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
	// 快速路径：已经初始化过，直接返回
	if runtimeE != nil {
		return runtimeE
	}

	// 慢路径：未初始化，使用 onceE.Do 保证只初始化一次
	onceE.Do(func() {
		runtimeE = engine.NewDefault()
	})
	return runtimeE
}
