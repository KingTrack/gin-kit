package safe

import (
	"context"
	"runtime"

	"github.com/KingTrack/gin-kit/kit/internal/tls/store"
	runtimee "github.com/KingTrack/gin-kit/kit/runtime"
)

func Go(ctx context.Context, fn func(ctx context.Context)) {
	defer func() {
		if rc := recover(); rc != nil {
			// 获取调用栈信息
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stack := string(buf[:n])

			runtimee.Get().LoggerRegistry().CrashLogger().Printf("goroutine panic recovered, stack:%s", stack)
		}
	}()

	// 从传入的 ctx 获取信息
	namespace := ctx.Value("namespace")
	traceID := ctx.Value("trace_id")

	go func() {
		if ns, ok := namespace.(string); ok {
			store.SetNamespace(ns)
		}
		if tid, ok := traceID.(string); ok {
			store.SetTraceID(tid)
		}

		// 直接执行函数，传递原始 ctx
		fn(ctx)
	}()
}
