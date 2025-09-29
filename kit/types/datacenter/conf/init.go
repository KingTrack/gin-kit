package conf

import (
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

type DiscoveryBackendEnum string

const (
	DiscoveryBackendConsul DiscoveryBackendEnum = "consul"
	DiscoveryBackendNacos  DiscoveryBackendEnum = "nacos"
)

type WatcherBackendEnum string

const (
	WatcherBackendConsul WatcherBackendEnum = "consul"
	WatcherBackendNacos  WatcherBackendEnum = "nacos"
	WatcherBackendEtcd   WatcherBackendEnum = "etcd"
)

type Config struct {
	DiscoveryBackendName DiscoveryBackendEnum `toml:"discovery_backend_name" json:"discovery_backend_name" yaml:"discovery_backend_name"`
	WatcherBackendName   WatcherBackendEnum   `toml:"watcher_backend_name" json:"watcher_backend_name" yaml:"watcher_backend_name"`
	Consul               Consul               `toml:"consul" json:"consul" yaml:"consul"`
	Nacos                Nacos                `toml:"nacos" json:"nacos" yaml:"nacos"`
	Etcd                 Etcd                 `toml:"etcd" json:"etcd" yaml:"etcd"`
}

type Consul struct {
	Addr               string `toml:"addr" json:"addr" yaml:"addr"`
	CheckIntervalSec   int64  `toml:"check_interval_sec" json:"check_interval_sec" yaml:"check_interval_sec"`
	CheckTimeoutSec    int64  `toml:"check_timeout_sec" json:"check_timeout_sec" yaml:"check_timeout_sec"`
	DeregisterAfterSec int64  `toml:"deregister_after_sec" json:"deregister_after_sec" yaml:"deregister_after_sec"`
}

type Nacos struct {
	Addresses   []string `toml:"address" json:"address" yaml:"address"`
	NamespaceID string   `toml:"namespace_id" json:"namespace_id" yaml:"namespace_id"`
	TimeoutMs   uint64   `toml:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
	CacheDir    string   `toml:"cache_dir" json:"cache_dir" yaml:"cache_dir"` // ./data/nacos/
	GroupName   string   `toml:"group_name" json:"group_name" yaml:"group_name"`
	ClusterName string   `toml:"cluster_name" json:"cluster_name" yaml:"cluster_name"`
}

func (n *Nacos) ToServerConfig() []constant.ServerConfig {
	if len(n.Addresses) == 0 {
		return nil
	}

	var configs []constant.ServerConfig
	for _, v := range n.Addresses {
		ipPort := strings.Split(v, ":")
		if len(ipPort) != 2 {
			return nil
		}
		port, err := strconv.ParseUint(ipPort[1], 10, 64)
		if err != nil {
			return nil
		}
		configs = append(configs, constant.ServerConfig{
			IpAddr: ipPort[0],
			Port:   port,
		})
	}

	return configs
}

func (n *Nacos) ToClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		NamespaceId:         n.NamespaceID,
		TimeoutMs:           n.TimeoutMs,
		CacheDir:            n.CacheDir,
		NotLoadCacheAtStart: true,
	}
}

type Etcd struct {
	Endpoints      []string `toml:"endpoints" json:"endpoints" yaml:"endpoints"`
	DialTimeoutSec int64    `toml:"dial_timeout_sec" json:"dial_timeout_sec" yaml:"dialTimeoutSec"`
}
