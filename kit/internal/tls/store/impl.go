package store

import (
	"sync"

	"github.com/modern-go/gls"
)

type Local struct {
	namespace string
	traceID   string
	peerName  string

	// 读写锁保护字段访问
	mu sync.RWMutex
}

// 全局协程存储
var (
	localStore = sync.Map{} // map[int64]*Local
)

// getGoID 获取当前 goroutine ID
func getGoID() int64 {
	return gls.GoID()
}

// getCurrentLocal 获取当前协程的本地存储
func getCurrentLocal() *Local {
	goid := getGoID()

	// 使用 LoadOrStore 确保原子性
	if local, loaded := localStore.LoadOrStore(goid, &Local{}); loaded {
		return local.(*Local)
	} else {
		// 新创建的对象已经在 LoadOrStore 中存储了
		return local.(*Local)
	}
}

func SetTraceID(traceID string) {
	local := getCurrentLocal()
	local.mu.Lock()
	defer local.mu.Unlock()
	local.traceID = traceID
}

func GetTraceID() string {
	local := getCurrentLocal()
	local.mu.RLock()
	defer local.mu.RUnlock()
	return local.traceID
}

func SetNamespace(namespace string) {
	local := getCurrentLocal()
	local.mu.Lock()
	defer local.mu.Unlock()
	local.namespace = namespace
}

func GetNamespace() string {
	local := getCurrentLocal()
	local.mu.RLock()
	defer local.mu.RUnlock()
	return local.namespace
}
