package middleware

import (
	"math/rand"
	"net/http"
	"time"

	clientcontext "github.com/KingTrack/gin-kit/kit/internal/httpclient/context"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
)

func Retry(config *conf.RetryerConfig) clientcontext.HandlerFunc {
	return func(cc *clientcontext.Context) {
		var (
			lastErr  error
			lastResp *http.Response
		)

		for attempt := 0; attempt <= config.RetryTimes; attempt++ {
			cc.Next() // 调用下一个中间件（通常是 call()）
			lastErr = cc.Err
			lastResp = cc.Resp.Response

			// 清理状态，为下次重试准备
			cc.Err = nil
			if cc.Resp.Response != nil {
				_ = cc.Resp.Response.Body.Close()
			}

			// 计算 backoff 时间
			sleepMs := calcBackoffMs(attempt, config.BaseDelayMs, config.MaxDelayMs, config.JitterFactor)
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		}

		cc.Err = lastErr
		cc.Resp.Response = lastResp
	}
}

func calcBackoffMs(attempt, baseDelayMs, maxDelayMs int, jitter float64) int {
	if attempt < 1 {
		attempt = 1
	}
	if baseDelayMs <= 0 {
		baseDelayMs = 1
	}
	if maxDelayMs < baseDelayMs {
		maxDelayMs = baseDelayMs
	}

	// 使用整数指数增长，限制最大延迟
	delayMs := baseDelayMs << (attempt - 1)
	if delayMs <= 0 || delayMs > maxDelayMs { // 防止移位溢出或超出上限
		delayMs = maxDelayMs
	}

	// jitter 限制在 [0, 1]
	if jitter < 0 {
		jitter = 0
	} else if jitter > 1 {
		jitter = 1
	}

	// ± jitter%
	if jitter > 0 {
		jitterRange := int(float64(delayMs) * jitter)
		if jitterRange > 0 {
			delayMs += rand.Intn(2*jitterRange+1) - jitterRange
			if delayMs < baseDelayMs {
				delayMs = baseDelayMs
			}
		}
	}

	return delayMs
}
