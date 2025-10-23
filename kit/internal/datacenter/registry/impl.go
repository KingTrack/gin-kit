package registry

import (
	"context"
	"sync"

	"github.com/KingTrack/gin-kit/kit/globals"
	"github.com/KingTrack/gin-kit/kit/internal/datacenter/balancer/define"
	"github.com/KingTrack/gin-kit/kit/internal/datacenter/balancer/roundroin"
	consulclient "github.com/KingTrack/gin-kit/kit/internal/datacenter/client/consul"
	etcdclient "github.com/KingTrack/gin-kit/kit/internal/datacenter/client/etcd"
	nacosclient "github.com/KingTrack/gin-kit/kit/internal/datacenter/client/nacos"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/discovery"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/watcher"
	httpclientconf "github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
	"github.com/pkg/errors"
)

type Registry struct {
	watcher   watcher.IWatcher
	discovery discovery.IDiscovery
	balancers map[string]define.IBalancer // 服务名字唯一维度
	mu        sync.RWMutex
}

func New() *Registry {
	return &Registry{
		balancers: make(map[string]define.IBalancer),
	}
}

func (r *Registry) Init(ctx context.Context, config *conf.Config) error {
	var consulClient *consulclient.Client
	var nacosClient *nacosclient.Client
	var etcdClient *etcdclient.Client

	switch config.DiscoveryBackendName {
	case conf.DiscoveryBackendConsul:
		if consulClient != nil {
			consulClient = consulclient.New()
			if err := consulClient.Init(ctx, &config.Consul); err != nil {
				return err
			}
		}
		r.discovery = consulClient
	case conf.DiscoveryBackendNacos:
		if nacosClient != nil {
			nacosClient = nacosclient.New()
			if err := nacosClient.Init(ctx, &config.Nacos); err != nil {
				return err
			}
		}
		r.discovery = nacosClient
	default:
		return errors.Errorf("unsupported dicovery backend: %s", config.DiscoveryBackendName)
	}

	switch config.WatcherBackendName {
	case conf.WatcherBackendConsul:
		if consulClient != nil {
			consulClient = consulclient.New()
			if err := consulClient.Init(ctx, &config.Consul); err != nil {
				return err
			}
		}
		r.watcher = consulClient
	case conf.WatcherBackendNacos:
		if nacosClient != nil {
			nacosClient = nacosclient.New()
			if err := nacosClient.Init(ctx, &config.Nacos); err != nil {
				return err
			}
		}
		r.watcher = nacosClient
	case conf.WatcherBackendEtcd:
		if etcdClient != nil {
			etcdClient = etcdclient.New()
			if err := etcdClient.Init(ctx, &config.Etcd); err != nil {
				return err
			}
		}
		r.watcher = etcdClient
	default:
		return errors.Errorf("unsupported watcher backend: %s", config.WatcherBackendName)
	}

	return nil
}

func (r *Registry) AddHTTPClient(ctx context.Context, config *httpclientconf.Config) error {
	if config.Discovery.IsDatacenter() {
		if err := r.AddServiceToDatacenter(ctx, config.ServiceName); err != nil {
			return errors.WithMessage(err, "datacenter registry add http client to datacenter failed")
		}
		return nil
	}

	if err := r.AddServiceToLocal(ctx, config); err != nil {
		return errors.WithMessage(err, "datacenter registry add http client to local failed")
	}
	return nil
}

func (r *Registry) AddServiceToDatacenter(ctx context.Context, serviceName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.balancers[serviceName]; ok {
		return nil
	}
	r.balancers[serviceName] = roundroin.New()

	go r.runDiscoveryWatcher(ctx, serviceName)

	return nil
}

func (r *Registry) AddServiceToLocal(ctx context.Context, config *httpclientconf.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.balancers[config.ServiceName]; ok {
		return nil
	}

	balancer := roundroin.New()
	balancer.Update(config.ToInstances())

	r.balancers[config.ServiceName] = balancer

	return nil
}

func (r *Registry) runDiscoveryWatcher(ctx context.Context, serviceName string) {
	for event := range r.discovery.WatchService(ctx, serviceName) {
		if err := event.Err; err != nil {
			globals.GetLogger().GenLogger().Printf("datacenter discovery watch %v instances error:%v", serviceName, err)
			continue
		}
		balancer := r.getBalancer(serviceName)
		if balancer == nil {
			globals.GetLogger().GenLogger().Printf("datacenter find %s balancer is nil", serviceName)
			continue
		}
		balancer.Update(event.Instances)
	}
}

func (r *Registry) PickInstance(serviceName string, meta map[string]string, skip *instance.Instance) (*instance.Instance, error) {
	balancer := r.getBalancer(serviceName)
	if balancer == nil {
		return nil, errors.Errorf("datacenter pick %s instance failed", serviceName)
	}
	return balancer.Pick(meta, skip)
}

func (r *Registry) getBalancer(serviceName string) define.IBalancer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.balancers[serviceName]
}

func (r *Registry) Register(instance *instance.Instance) error {
	return r.discovery.Register(instance)
}

func (r *Registry) Deregister(instance *instance.Instance) error {
	return r.discovery.Deregister(instance)
}

func (r *Registry) WatchKV(ctx context.Context, key string) <-chan watcher.Event {
	return r.watcher.WatchKV(ctx, key)
}

func (r *Registry) WatchPrefix(ctx context.Context, prefix string) <-chan watcher.Event {
	return r.watcher.WatchPrefix(ctx, prefix)
}
