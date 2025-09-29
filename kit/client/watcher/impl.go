package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/watcher"
)

func WatchKV(ctx context.Context, key string) <-chan watcher.Event {
	return runtime.Get().DatacenterRegistry().WatchKV(ctx, key)
}
