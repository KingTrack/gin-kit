package client

import (
	"context"
	"fmt"
	"time"

	"github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/discovery"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/watcher"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Client struct {
	config *conf.Consul
	client *api.Client
}

func New() *Client {
	return &Client{}
}

func (c *Client) Init(ctx context.Context, config *conf.Consul) error {
	c.config = config

	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Addr

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return errors.WithMessage(err, "consul client create failed")
	}

	c.client = client

	return nil
}

func (c *Client) Register(instance *instance.Instance) error {
	check := &api.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%d", instance.IP, instance.Port),
		Interval:                       fmt.Sprintf("%ds", c.config.CheckIntervalSec),
		Timeout:                        fmt.Sprintf("%ds", c.config.CheckTimeoutSec),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.config.DeregisterAfterSec),
	}

	reg := &api.AgentServiceRegistration{
		ID:      instance.ServiceID(),
		Name:    instance.ServiceName,
		Address: instance.IP,
		Port:    instance.Port,
		Meta:    instance.GetMeta(),
		Check:   check,
	}

	if err := c.client.Agent().ServiceRegister(reg); err != nil {
		return errors.WithMessage(err, "consul client register service to consul failed")
	}

	return nil
}

func (c *Client) Deregister(instance *instance.Instance) error {
	if err := c.client.Agent().ServiceDeregister(instance.ServiceID()); err != nil {
		return errors.WithMessage(err, "consul client deregister service from consul failed")
	}

	return nil
}

func (c *Client) WatchService(ctx context.Context, serviceName string) <-chan discovery.Event {
	out := make(chan discovery.Event, 10)

	go func() {
		defer close(out)

		health := c.client.Health()
		var lastIndex uint64
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			serviceEntries, meta, err := health.Service(serviceName, "", true, &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  time.Duration(60) * time.Second,
			})
			if err != nil {
				select {
				case out <- discovery.Event{Err: err, Instances: nil}:
				case <-ctx.Done():
					return
				}
				continue
			}

			if meta.LastIndex == lastIndex {
				continue
			}
			lastIndex = meta.LastIndex

			if len(serviceEntries) == 0 {
				continue
			}

			var instances []instance.Instance
			for _, v := range serviceEntries {
				instances = append(instances, instance.Instance{
					ServiceName: v.Service.Service,
					IP:          v.Service.Address,
					Port:        v.Service.Port,
					Weight:      instance.GetWeight(v.Service.Meta),
					Meta:        instance.RebuildMeta(v.Service.Meta),
				})
			}
			select {
			case out <- discovery.Event{Err: nil, Instances: instances}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (c *Client) WatchKV(ctx context.Context, key string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)
	go func() {
		defer close(out)

		kv := c.client.KV()
		var lastIndex uint64
		for {
			pair, meta, err := kv.Get(key, &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  time.Duration(60) * time.Second,
			})
			if err != nil {
				select {
				case out <- watcher.Event{Err: err, Data: nil}:
				case <-ctx.Done():
					return
				}
				continue
			}

			if meta.LastIndex == lastIndex {
				continue
			}
			lastIndex = meta.LastIndex

			if pair == nil {
				continue
			}

			select {
			case out <- watcher.Event{Data: map[string][]byte{key: pair.Value}}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (c *Client) WatchPrefix(ctx context.Context, prefix string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)
	go func() {
		defer close(out)

		kv := c.client.KV()
		var lastIndex uint64
		for {
			pairs, meta, err := kv.List(prefix, &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  time.Duration(60) * time.Second,
			})
			if err != nil {
				select {
				case out <- watcher.Event{Err: err, Data: nil}:
				case <-ctx.Done():
					return
				}
				continue
			}

			if meta.LastIndex == lastIndex {
				continue
			}
			lastIndex = meta.LastIndex

			if len(pairs) == 0 {
				continue
			}

			data := make(map[string][]byte, len(pairs))
			for _, p := range pairs {
				data[p.Key] = p.Value
			}

			select {
			case out <- watcher.Event{Err: nil, Data: data}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
