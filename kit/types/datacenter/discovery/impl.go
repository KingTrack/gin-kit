package discovery

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
)

type Event struct {
	Err       error
	Instances []instance.Instance
}

type IDiscovery interface {
	Register(instance *instance.Instance) error
	Deregister(instance *instance.Instance) error
	WatchService(ctx context.Context, serviceName string) <-chan Event
}
