package messagebroker

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

type PageRedisMessageBroker struct {
	client    *redis.Client
	publisher *PageMessageBrokerPublisher
}

func NewPageRedisMessageBroker(producer *redis.Client) *PageRedisMessageBroker {
	return &PageRedisMessageBroker{
		client:    producer,
		publisher: NewPageMessageBrokerPublisher(producer),
	}
}

func (m *PageRedisMessageBroker) Publisher() *PageMessageBrokerPublisher {
	return m.publisher
}

func (m *PageRedisMessageBroker) Close() error {
	if err := m.client.Close(); err != nil {
		return errors.New("redis message broker error while closing")
	}

	return nil
}
