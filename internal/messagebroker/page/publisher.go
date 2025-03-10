package messagebroker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/Stuhub-io/core/domain"
	"github.com/redis/go-redis/v9"
)

type PageMessageBrokerPublisher struct {
	client *redis.Client
}

func NewPageMessageBrokerPublisher(client *redis.Client) *PageMessageBrokerPublisher {
	return &PageMessageBrokerPublisher{
		client: client,
	}
}

func (p *PageMessageBrokerPublisher) Created(ctx context.Context, page *domain.Page) error {
	return p.publish(ctx, "page.created", page)
}

func (p *PageMessageBrokerPublisher) Deleted(ctx context.Context, id string) error {
	return p.publish(ctx, "page.deleted", id)
}

func (p *PageMessageBrokerPublisher) Updated(ctx context.Context, page *domain.Page) error {
	return p.publish(ctx, "page.updated", page)
}

func (p *PageMessageBrokerPublisher) publish(ctx context.Context, channel string, event any) error {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(event); err != nil {
		return errors.New("failed to encode event")
	}

	res := p.client.Publish(ctx, channel, b.Bytes())
	if err := res.Err(); err != nil {
		return errors.New("failed to publish event")
	}

	return nil
}
