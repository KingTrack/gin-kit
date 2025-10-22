package conf

import "github.com/Shopify/sarama"

type Config struct {
	Name         string              `toml:"name" json:"name" yaml:"name"`
	Addrs        []string            `toml:"addrs" json:"addrs" yaml:"addrs"`
	RequiredAcks sarama.RequiredAcks `toml:"required_acks" json:"required_acks" yaml:"required_acks"`
	RetryMax     int                 `toml:"retry_max" json:"retry_max" yaml:"retry_max"`
}
