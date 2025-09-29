package etcd

import (
	"context"
	"time"

	"github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	"github.com/KingTrack/gin-kit/kit/types/datacenter/watcher"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	config *conf.Etcd
	client *clientv3.Client
}

func New() *Client {
	return &Client{}
}

func (c *Client) Init(ctx context.Context, config *conf.Etcd) error {
	c.config = config

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: time.Duration(config.DialTimeoutSec) * time.Second,
	})
	if err != nil {
		return errors.WithMessage(err, "etcd client create failed")
	}
	c.client = client

	return nil
}

func (c *Client) WatchKV(ctx context.Context, key string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)

	go func() {
		defer close(out)

		currentResp, err := c.client.Get(ctx, key)
		if err != nil {
			select {
			case out <- watcher.Event{Err: err, Data: nil}: // 只返回错误，不退出监听
			case <-ctx.Done():
				return
			}
		}
		if len(currentResp.Kvs) > 0 {
			for _, v := range currentResp.Kvs {
				if string(v.Key) == key {
					select {
					case out <- watcher.Event{Err: nil, Data: map[string][]byte{key: v.Value}}:
					case <-ctx.Done():
						return
					}
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			watchChan := c.client.Watch(ctx, key)
			for {
				select {
				case watchResp := <-watchChan:
					if watchResp.Canceled {
						select {
						case out <- watcher.Event{Err: rpctypes.ErrCompacted, Data: nil}:
						case <-ctx.Done():
							return
						}

						select {
						case <-time.After(time.Second * 1):
						case <-ctx.Done():
							return
						}

						break // 跳出循环进行重新连接
					}

					if len(watchResp.Events) == 0 {
						continue
					}

					for _, v := range watchResp.Events {
						switch v.Type {
						case mvccpb.PUT:
							if string(v.Kv.Key) == key {
								select {
								case out <- watcher.Event{Err: nil, Data: map[string][]byte{key: v.Kv.Value}}:
								case <-ctx.Done():
									return
								}
							}
						}
					}
				}
			}

		}
	}()

	return out
}

func (c *Client) WatchPrefix(ctx context.Context, prefix string) <-chan watcher.Event {
	out := make(chan watcher.Event, 1)

	go func() {
		defer close(out)

		currentResp, err := c.client.Get(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			select {
			case out <- watcher.Event{Err: err, Data: nil}: // 只返回错误，不退出监听
			case <-ctx.Done():
				return
			}
		}

		if len(currentResp.Kvs) > 0 {
			data := make(map[string][]byte, len(currentResp.Kvs))
			for _, v := range currentResp.Kvs {
				data[string(v.Key)] = v.Value
			}
			select {
			case out <- watcher.Event{Err: nil, Data: data}:
			case <-ctx.Done():
				return
			}
		}

		for {
			select {
			case <-ctx.Done():
			default:
			}

			watchChan := c.client.Watch(ctx, prefix, clientv3.WithPrefix())
			for {
				select {
				case watchResp := <-watchChan:
					if watchResp.Canceled {
						select {
						case out <- watcher.Event{Err: rpctypes.ErrCompacted, Data: nil}:
						case <-ctx.Done():
							return
						}

						select {
						case <-time.After(time.Second * 1):
						case <-ctx.Done():
							return
						}

						break
					}

					if len(watchResp.Events) == 0 {
						continue
					}

					data := make(map[string][]byte, len(watchResp.Events))
					for _, v := range watchResp.Events {
						switch v.Type {
						case mvccpb.PUT:
							data[string(v.Kv.Key)] = v.Kv.Value
						}
					}
					if len(data) == 0 {
						continue
					}

					select {
					case out <- watcher.Event{Err: nil, Data: data}:
					case <-ctx.Done():
						return
					}

				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}
