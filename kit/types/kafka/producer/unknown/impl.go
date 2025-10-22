package unknown

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	ErrUnknownProducer = errors.New("unknown kafka producer")
)

type Producer struct{}

func New() sarama.SyncProducer {
	return &Producer{}
}

func (p *Producer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, ErrUnknownProducer
}

func (p *Producer) SendMessages(msgs []*sarama.ProducerMessage) error {
	return ErrUnknownProducer
}

func (p *Producer) Close() error {
	return ErrUnknownProducer
}

func (p *Producer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return sarama.ProducerTxnFlagFatalError
}

func (p *Producer) IsTransactional() bool {
	return false
}

func (p *Producer) BeginTxn() error {
	return ErrUnknownProducer
}

func (p *Producer) CommitTxn() error {
	return ErrUnknownProducer
}

func (p *Producer) AbortTxn() error {
	return ErrUnknownProducer
}

func (p *Producer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return ErrUnknownProducer
}

func (p *Producer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	return ErrUnknownProducer
}
