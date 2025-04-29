package build

import (
	"github.com/IBM/sarama"
	"github.com/cockroachdb/errors"

	"github.com/powerdigital/go-micro/pkg/producer"
)

func (b *Builder) Producer(brokers []string, topic string) (*producer.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "producer creating")
	}

	return &producer.Producer{
		SyncProducer:    syncProducer,
		CreateUserTopic: topic,
	}, nil
}
