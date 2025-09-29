package middleware

import (
	"strconv"
	"strings"
	"sync"
	"time"

	enginetls "github.com/KingTrack/gin-kit/kit/internal/tls/store"

	"github.com/KingTrack/gin-kit/kit/types/httpserver/middleware/logger"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	config      logger.IConfig
	messagePool sync.Pool
}

func New(config logger.IConfig) *Middleware {
	return &Middleware{
		config: config,
		messagePool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

func (m *Middleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.config == nil {
			c.Next()
			return
		}

		// 执行业务逻辑
		c.Next()

		// 记录访问日志
		runtime.Get().LoggerRegistry().AccessLogger().Print(m.getMessage(c))
	}
}

func (m *Middleware) getMessage(c *gin.Context) string {
	// 获取上下文信息
	namespace := enginetls.GetNamespace()
	traceID := enginetls.GetTraceID()

	// 默认值
	startTime := time.Now()
	var appCode string
	var peerName string
	var clientIP string
	var requestBody []byte
	var responseBody []byte

	// 从 Context Registry 获取数据（已经根据配置进行了大小限制和存储控制）
	if cc := runtime.Get().ContextRegistry().Load(c); cc != nil {
		startTime = cc.GetRequestStartTime()
		appCode = cc.GetAppCode()
		peerName = cc.GetPeerName()
		clientIP = cc.GetClientIP()
		requestBody = cc.GetRequestBody()   // 已经根据配置处理
		responseBody = cc.GetResponseBody() // 已经根据配置处理
	}

	// 从 pool 获取 builder
	builder := m.messagePool.Get().(*strings.Builder)
	builder.Reset()

	tryAppend := func(key, value string) {
		if builder.Len() > 0 {
			builder.WriteString("|")
		}
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(value)
	}

	// 请求耗时计算
	duration := time.Since(startTime).Microseconds()

	// 按顺序添加基础字段
	tryAppend("time", time.Now().Format("2006-01-02 15:04:05.000000"))
	tryAppend(m.config.GetRequestStartTimeKey(), startTime.Format("2006-01-02 15:04:05.000000"))
	tryAppend(m.config.GetRequestDurationMsKey(), strconv.FormatInt(duration, 10))
	tryAppend("trace_id", traceID)
	tryAppend("peer_name", peerName)
	tryAppend(m.config.GetClientIPKey(), clientIP)
	tryAppend("namespace", namespace)
	tryAppend("req_method", c.Request.Method)
	tryAppend(m.config.GetHTTPStatusKey(), strconv.Itoa(c.Writer.Status()))
	tryAppend("req_uri", c.Request.RequestURI)
	tryAppend("method_name", m.config.GetMethodNameValue(c.Request.URL.Path))
	tryAppend(m.config.GetAppCodeKey(), appCode)
	tryAppend("req_body", string(requestBody))
	tryAppend("resp_body", string(responseBody))

	// 获取结果并归还
	result := builder.String()
	builder.Reset()
	m.messagePool.Put(builder)

	return result
}
