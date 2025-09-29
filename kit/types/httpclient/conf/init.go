package conf

import (
	"strconv"
	"strings"

	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
)

type DiscoveryEnum string

func (e DiscoveryEnum) IsDatacenter() bool {
	return e == DiscoveryDatacenter
}

const (
	DiscoveryDatacenter DiscoveryEnum = "datacenter"
)

type Config struct {
	ServiceName string        `toml:"service_name" json:"service_name" yaml:"service_name"`
	Endpoints   []string      `toml:"endpoints" json:"endpoints" yaml:"endpoints"`
	Discovery   DiscoveryEnum `toml:"discovery" json:"discovery" yaml:"discovery"`
}

func (c *Config) ToInstances() []instance.Instance {
	var instances []instance.Instance
	for _, v := range c.Endpoints {
		ipPort := strings.Split(v, ":")
		if len(ipPort) != 2 {
			continue
		}
		ip := ipPort[0]
		port, err := strconv.Atoi(ipPort[1])
		if err != nil {
			continue
		}
		instances = append(instances, instance.Instance{
			ServiceName: c.ServiceName,
			IP:          ip,
			Port:        port,
			Weight:      100,
			Meta:        nil,
		})
	}
	return instances
}
