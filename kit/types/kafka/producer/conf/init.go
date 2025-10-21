package conf

import "github.com/Shopify/sarama"

type Config struct {
	Name            string              `toml:"name" json:"name" yaml:"name"`
	Addrs           []string            `toml:"addrs" json:"addrs" yaml:"addrs"`
	RequiredAcks    sarama.RequiredAcks `toml:"required_acks" json:"required_acks" yaml:"required_acks"`
	RetryMax        int                 `toml:"retry_max" json:"retry_max" yaml:"retry_max"`
	ReturnSuccesses bool                `toml:"return_successes" json:"return_successes" yaml:"return_successes"` // 每条发送成功的消息都会返回 Partition 和 Offset，方便做日志或指标统计
}
