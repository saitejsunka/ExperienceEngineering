package producers

import (
	"context"

	"dbexp/backend/service/consumers"
)

// HealthProducer defines the interface for formulating health check responses.
type HealthProducer interface {
	CheckHealth(ctx context.Context) (int32, string)
}

type healthProducer struct {
	consumer consumers.HealthConsumer
}

// NewHealthProducer creates a new HealthProducer instance, injecting the HealthConsumer dependency.
func NewHealthProducer(c consumers.HealthConsumer) HealthProducer {
	return &healthProducer{
		consumer: c,
	}
}

// CheckHealth processes the health check request and forwards it to the consumer.
func (p *healthProducer) CheckHealth(ctx context.Context) (int32, string) {
	err := p.consumer.PingDatabase(ctx)
	if err != nil {
		return 500, "Database is unreachable"
	}

	return 200, "Service is healthy and Database is reachable"
}
