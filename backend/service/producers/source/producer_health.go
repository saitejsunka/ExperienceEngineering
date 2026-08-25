package source

import (
	"context"
	"fmt"

	consumersource "dbexp/backend/service/consumers/source"
)

// HealthProducer defines the interface for health business logic.
type HealthProducer interface {
	CheckHealth(ctx context.Context) (int32, string, error)
}

// HealthProducerImpl implements the HealthProducer interface.
type HealthProducerImpl struct {
	Consumer consumersource.HealthConsumer
}

// CheckHealth contains the business logic to verify system health.
// It delegates to the Consumer to verify database connectivity.
func (p *HealthProducerImpl) CheckHealth(ctx context.Context) (int32, string, error) {
	err := p.Consumer.PingDatabase(ctx)
	if err != nil {
		return 500, fmt.Sprintf("Service is healthy, but Database is unreachable: %v", err), nil
	}

	return 200, "Service is healthy and Database is reachable", nil
}
