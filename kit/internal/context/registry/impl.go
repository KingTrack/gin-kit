package registry

import (
	"fmt"
	"sync"

	corecontext "github.com/KingTrack/gin-kit/kit/internal/context/core"
	"github.com/gin-gonic/gin"
)

type Registry struct {
	store sync.Map
	pool  *corecontext.Pool
}

func New() *Registry {
	return &Registry{
		pool: corecontext.NewPool(),
	}
}

func (r *Registry) GetPool() *corecontext.Pool {
	return r.pool
}

func (r *Registry) Store(c *gin.Context, cc *corecontext.Context) {
	key := getKey(c)
	r.store.Store(key, cc)
}

func (r *Registry) Load(c *gin.Context) *corecontext.Context {
	key := getKey(c)
	if value, ok := r.store.Load(key); ok {
		if ctx, ok := value.(*corecontext.Context); ok {
			return ctx
		}
	}
	return nil
}

func (r *Registry) Remove(c *gin.Context) {
	key := getKey(c)
	if value, ok := r.store.Load(key); ok {
		if cc, ok := value.(*corecontext.Context); ok {
			cc.Cleanup()
			r.pool.Put(cc)
		}
		r.store.Delete(key)
	}
}

func getKey(c *gin.Context) string {
	return fmt.Sprintf("%p", c)
}
