package conf

import (
	"net"
	"net/url"
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
	TimeoutMs   int64         `toml:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
}

func (c *Config) ToInstances() []instance.Instance {
	var instances []instance.Instance
	for _, v := range c.Endpoints {
		if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
			u, err := url.Parse(v)
			if err != nil {
				continue
			}

			ip, portStr, err := net.SplitHostPort(u.Host)
			if err != nil {
				continue
			}
			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}

			instances = append(instances, instance.Instance{
				ServiceName: c.ServiceName,
				Schema:      u.Scheme,
				IP:          ip,
				Port:        port,
				Weight:      100,
				Meta:        nil,
			})
			continue
		}

		ip, portStr, err := net.SplitHostPort(v)
		if err != nil {
			continue
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		instances = append(instances, instance.Instance{
			ServiceName: c.ServiceName,
			Schema:      "http",
			IP:          ip,
			Port:        port,
			Weight:      100,
			Meta:        nil,
		})
	}
	return instances
}
