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
	ServiceName        string        `toml:"service_name" json:"service_name" yaml:"service_name"`
	Endpoints          []string      `toml:"endpoints" json:"endpoints" yaml:"endpoints"`
	Discovery          DiscoveryEnum `toml:"discovery" json:"discovery" yaml:"discovery"`
	MaxIdleConns       int           `toml:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	IdleConnTimeoutSec int           `toml:"idle_conn_timeout_sec" json:"idle_conn_timeout_sec" yaml:"idle_conn_timeout_sec"`
	TimeoutMs          int           `toml:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
	ProxyURL           string        `toml:"proxy_url" json:"proxy_url" yaml:"proxy_url"`
	RetryerConfig      RetryerConfig `toml:"retryer_config" json:"retryer_config" yaml:"retryer_config"`
}

type RetryerConfig struct {
	RetryTimes   int     `toml:"retry_times" json:"retry_times" yaml:"retry_times"`
	BaseDelayMs  int     `toml:"base_delay_ms" json:"base_delay_ms" yaml:"base_delay_ms"`
	MaxDelayMs   int     `toml:"max_delay_ms" json:"max_delay_ms" yaml:"max_delay_ms"`
	JitterFactor float64 `toml:"jitter_factor" json:"jitter_factor" yaml:"jitter_factor"`
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
