package producers

import (
	consumer "dbexp/backend/service/consumers/source"
	"dbexp/backend/service/producers/source"
)

// InitializeHealthProducer creates a producer, injecting the necessary consumer.
func InitializeHealthProducer(c consumer.HealthConsumer) source.HealthProducer {
	return &source.HealthProducerImpl{
		Consumer: c,
	}
}
