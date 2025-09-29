package watcher

import "context"

type Event struct {
	Err  error
	Data map[string][]byte
}

type IWatcher interface {
	WatchKV(ctx context.Context, key string) <-chan Event
	WatchPrefix(ctx context.Context, prefix string) <-chan Event
}
