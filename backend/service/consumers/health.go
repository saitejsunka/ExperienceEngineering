package consumers

import (
	"context"

	"dbexp/backend/service/configs"
	"dbexp/backend/service/observability"
)

// HealthConsumer defines the interface for executing health queries against the database.
type HealthConsumer interface {
	PingDatabase(ctx context.Context) error
}

// mockHealthConsumer implements the HealthConsumer interface with mocked database calls.
type mockHealthConsumer struct {
	config *configs.AppConfig
	telemetry observability.Telemetry
}

// NewMockHealthConsumer creates a mock consumer injecting config and telemetry.
func NewMockHealthConsumer(cfg *configs.AppConfig, tel observability.Telemetry) HealthConsumer {
	return &mockHealthConsumer{
		config:    cfg,
		telemetry: tel,
	}
}

// PingDatabase simulates checking the database health.
func (m *mockHealthConsumer) PingDatabase(ctx context.Context) error {
	// Log when the request came in
	m.telemetry.LogInfo("Received PingDatabase request", map[string]interface{}{
		"target_db": m.config.DatabaseHost,
		"db_name":   m.config.DatabaseName,
	})

	// STOPPED HERE: We are stopping here as per instructions.
	// In the future, the actual Cloud SQL connection will be pinged here.
	
	// Log the response
	m.telemetry.LogInfo("PingDatabase response", map[string]interface{}{
		"status": 200,
		"message": "Mocking a successful database connection",
	})
	
	return nil // Mocking a successful database connection
}
