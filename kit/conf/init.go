package conf

import (
	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/plugin/decoder"
	"github.com/KingTrack/gin-kit/kit/plugin/source"
	datacenterconfig "github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	serverconfig "github.com/KingTrack/gin-kit/kit/types/httpserver/conf"
	loggerconfig "github.com/KingTrack/gin-kit/kit/types/logger/conf"
	metricconfig "github.com/KingTrack/gin-kit/kit/types/metric/conf"
	mysqlconfig "github.com/KingTrack/gin-kit/kit/types/mysql/conf"
	redisconfig "github.com/KingTrack/gin-kit/kit/types/redis/conf"
	tracerconfig "github.com/KingTrack/gin-kit/kit/types/tracer/conf"
)

type Config struct {
	Namespace  tlscontext.Value        `toml:"namespace" json:"namespace" yaml:"namespace"`
	Hostname   string                  `toml:"hostname" json:"hostname" yaml:"hostname"`
	Server     serverconfig.Config     `toml:"httpserver" json:"httpserver" yaml:"httpserver"`
	Metric     metricconfig.Config     `toml:"metric" json:"metric" yaml:"metric"`
	Tracer     tracerconfig.Config     `toml:"tracer" json:"tracer" yaml:"tracer"`
	Logger     loggerconfig.Config     `toml:"logger" json:"logger" yaml:"logger"`
	Datacenter datacenterconfig.Config `toml:"datacenter" json:"datacenter" yaml:"datacenter"`
	MySQL      []mysqlconfig.Config    `toml:"mysql" json:"mysql" yaml:"mysql"`
	Redis      []redisconfig.Config    `toml:"redis" json:"redis" yaml:"redis"`
}

type Namespace struct {
	RootPath string
	Source   source.ISource
	Decoder  decoder.IDecoder
}
