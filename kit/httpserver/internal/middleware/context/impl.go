package middleware

import (
	"time"

	corecontext "github.com/KingTrack/gin-kit/kit/internal/context/core"
	enginetls "github.com/KingTrack/gin-kit/kit/internal/tls/store"
	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Config IContextConfig
}

func (m *Middleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		if runtime.Get() == nil {
			return
		}

		// 1. 创建并初始化 context
		startTime := time.Now()
		newContext := runtime.Get().ContextRegistry().GetPool().Get()
		newContext.SetRequestStartTime(startTime)

		// 2. 存储到管理器
		runtime.Get().ContextRegistry().Store(c, newContext)

		// 3. 解析上下文信息
		if m.Config != nil {
			m.parseContext(c, newContext)
		}

		// 4. 执行后续逻辑
		c.Next()

		// 5. 清理 context
		runtime.Get().ContextRegistry().Remove(c)
	}
}

func (m *Middleware) parseContext(c *gin.Context, cc *corecontext.Context) {
	namespace := m.Config.ParseNamespace(c.Request.Header)
	clientIP := m.Config.ParseClientIP(c.Request.Header, c.Request.RequestURI)

	if len(namespace) > 0 {
		enginetls.SetNamespace(namespace)
	}
	if len(clientIP) > 0 {
		cc.SetClientIP(clientIP)
	}
}
