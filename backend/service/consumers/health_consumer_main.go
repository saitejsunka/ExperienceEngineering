package consumers

import (
	"dbexp/backend/service/configs"
	"dbexp/backend/service/consumers/source"
	"dbexp/backend/service/observability"
)

// InitializeHealthConsumer creates a consumer injecting config and telemetry.
func InitializeHealthConsumer(cfg *configs.AppConfig, tel observability.Telemetry) source.HealthConsumer {
	return &source.HealthConsumerImpl{
		Config:    cfg,
		Telemetry: tel,
	}
}
