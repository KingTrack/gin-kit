package conf

import (
	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/plugin/decoder"
	"github.com/KingTrack/gin-kit/kit/plugin/source"
	datacenterconf "github.com/KingTrack/gin-kit/kit/types/datacenter/conf"
	httpclientconf "github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
	serverconf "github.com/KingTrack/gin-kit/kit/types/httpserver/conf"
	kafkaconsumerconf "github.com/KingTrack/gin-kit/kit/types/kafka/consumer/conf"
	kafkaproducerconf "github.com/KingTrack/gin-kit/kit/types/kafka/producer/conf"
	loggerconf "github.com/KingTrack/gin-kit/kit/types/logger/conf"
	metricconf "github.com/KingTrack/gin-kit/kit/types/metric/conf"
	mysqlconf "github.com/KingTrack/gin-kit/kit/types/mysql/conf"
	redisconf "github.com/KingTrack/gin-kit/kit/types/redis/conf"
	tracerconf "github.com/KingTrack/gin-kit/kit/types/tracer/conf"
)

type Config struct {
	Namespace     tlscontext.Value           `toml:"namespace" json:"namespace" yaml:"namespace"`
	Hostname      string                     `toml:"hostname" json:"hostname" yaml:"hostname"`
	Server        serverconf.Config          `toml:"httpserver" json:"httpserver" yaml:"httpserver"`
	Metric        metricconf.Config          `toml:"metric" json:"metric" yaml:"metric"`
	Tracer        tracerconf.Config          `toml:"tracer" json:"tracer" yaml:"tracer"`
	Logger        loggerconf.Config          `toml:"logger" json:"logger" yaml:"logger"`
	Datacenter    datacenterconf.Config      `toml:"datacenter" json:"datacenter" yaml:"datacenter"`
	HTTPClient    []httpclientconf.Config    `toml:"http_client" json:"http_client" yaml:"http_client"`
	MySQL         []mysqlconf.Config         `toml:"mysql" json:"mysql" yaml:"mysql"`
	Redis         []redisconf.Config         `toml:"redis" json:"redis" yaml:"redis"`
	KafkaProducer []kafkaproducerconf.Config `toml:"kafka_producer" json:"kafka_producer" yaml:"kafka_producer"`
	KafkaConsumer []kafkaconsumerconf.Config `toml:"kafka_consumer" json:"kafka_consumer" yaml:"kafka_consumer"`
}

type Namespace struct {
	RootPath string
	Source   source.ISource
	Decoder  decoder.IDecoder
}
