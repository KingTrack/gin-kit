package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/globals"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/discovery"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/watcher"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
)

type Client struct {
	config       *conf.Nacos
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
}

func New() *Client {
	return &Client{}
}

func (c *Client) Init(ctx context.Context, config *conf.Nacos) error {
	c.config = config

	clientConfig := config.ToClientConfig()
	serverConfigs := config.ToServerConfig()

	if clientConfig == nil {
		return errors.New("nacos client config is invalid")
	}
	if len(serverConfigs) == 0 {
		return errors.New("nacos httpserver configs is invalid")
	}

	clientParams := vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	}

	namingClient, err := clients.NewNamingClient(clientParams)
	if err != nil {
		return errors.WithMessage(err, "nacos client create naming client failed")
	}

	configClient, err := clients.NewConfigClient(clientParams)
	if err != nil {
		return errors.WithMessage(err, "nacos client create config client failed")
	}

	logger.SetLogger(globals.GetLogger().GenLogger())

	c.namingClient = namingClient
	c.configClient = configClient
	return nil
}

// Ephemeral:   true
// 自动心跳: 客户端会定期向 Nacos 服务器发送心跳
// 自动故障检测: 如果心跳停止（通常 15-30 秒无心跳），实例会被自动标记为不健康/删除
// 服务发现: 其他服务查询时，只会返回健康的实例
func (c *Client) Register(instance *instance.Instance) error {
	ok, err := c.namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          instance.IP,
		Port:        uint64(instance.Port),
		Weight:      100,
		Enable:      true,
		Healthy:     true,
		Metadata:    instance.GetMeta(),
		ServiceName: instance.ServiceName,
		ClusterName: c.config.ClusterName,
		GroupName:   c.config.GroupName,
		Ephemeral:   true, // 自动心跳: 客户端会定期向 Nacos 服务器发送心跳
	})
	if err != nil {
		return errors.WithMessage(err, "nacos client register service to nacos failed")
	}

	if ok {
		return nil
	}

	return errors.New("nacos client register service to nacos failed")
}

func (c *Client) Deregister(instance *instance.Instance) error {
	ok, err := c.namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          instance.IP,
		Port:        uint64(instance.Port),
		Cluster:     c.config.ClusterName,
		ServiceName: instance.ServiceName,
		GroupName:   c.config.GroupName,
		Ephemeral:   true, // 与注册时保持一致
	})
	if err != nil {
		return errors.WithMessage(err, "nacos client deregister service from nacos failed")
	}

	if ok {
		return nil
	}

	return errors.New("nacos client deregister service from nacos failed")
}

func (c *Client) WatchService(ctx context.Context, serviceName string) <-chan discovery.Event {
	out := make(chan discovery.Event, 10)
	go func() {
		defer close(out)

		err := c.namingClient.Subscribe(&vo.SubscribeParam{
			ServiceName: serviceName,
			Clusters:    []string{c.config.ClusterName},
			GroupName:   c.config.GroupName,
			SubscribeCallback: func(all []model.Instance, err error) {
				// 检查context是否已取消
				select {
				case <-ctx.Done():
					return // context已取消，不再处理事件
				default:
				}

				if err != nil {
					select {
					case out <- discovery.Event{Err: err, Instances: nil}:
					case <-ctx.Done():
					}
					return
				}

				var instances []instance.Instance
				for _, v := range all {
					if !v.Healthy {
						continue
					}
					instances = append(instances, instance.Instance{
						ServiceName: v.ServiceName,
						IP:          v.Ip,
						Port:        int(v.Port),
						Weight:      instance.GetWeight(v.Metadata),
						Meta:        instance.RebuildMeta(v.Metadata),
					})
				}
				// 发送事件时也要检查context
				select {
				case out <- discovery.Event{Err: nil, Instances: instances}:
				case <-ctx.Done():
					return
				}
			},
		})
		if err != nil {
			select {
			case <-ctx.Done():
			case out <- discovery.Event{Err: err, Instances: nil}:
			}
			return
		}

		// 等待取消
		<-ctx.Done()

		if err := c.namingClient.Unsubscribe(&vo.SubscribeParam{
			ServiceName: serviceName,
			Clusters:    []string{c.config.ClusterName},
			GroupName:   c.config.GroupName,
		}); err != nil {
			select {
			case out <- discovery.Event{Err: err, Instances: nil}:
			case <-ctx.Done():
			}
			return
		}
	}()

	return out
}

func (c *Client) WatchKV(ctx context.Context, key string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)

	go func() {
		defer close(out)

		// ListenConfig是异步的，立即返回
		if err := c.configClient.ListenConfig(vo.ConfigParam{
			DataId: key,
			Group:  c.config.GroupName,
			OnChange: func(_, _, _, data string) {
				select {
				case out <- watcher.Event{Err: nil, Data: map[string][]byte{key: []byte(data)}}:
				case <-ctx.Done():
					return
				}
			},
		}); err != nil {
			select {
			case out <- watcher.Event{Err: err, Data: nil}:
			case <-ctx.Done():
			}
			return
		}

		// 等待context取消
		<-ctx.Done()

		if err := c.configClient.CancelListenConfig(vo.ConfigParam{
			DataId: key,
			Group:  c.config.GroupName,
		}); err != nil {
			select {
			case out <- watcher.Event{Err: err, Data: nil}:
			case <-ctx.Done():
			}
			return
		}
	}()

	return out
}

func (c *Client) WatchPrefix(ctx context.Context, prefix string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)
	go func() {
		defer close(out)

		select {
		case out <- watcher.Event{Err: errors.New("WatchPrefix is not supported in Nacos SDK, consider using WatchKV for specific keys")}:
		case <-ctx.Done():
			return
		}

		// 等待取消
		<-ctx.Done()
	}()
	return out
}
