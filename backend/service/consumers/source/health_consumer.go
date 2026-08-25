package source

import (
	"context"
	"database/sql"
	"fmt"

	"dbexp/backend/service/configs"
	"dbexp/backend/service/observability"
	_ "github.com/lib/pq"
)

// HealthConsumer defines the interface for executing health queries against the database.
type HealthConsumer interface {
	PingDatabase(ctx context.Context) error
}

// HealthConsumerImpl implements the HealthConsumer interface with actual database calls.
type HealthConsumerImpl struct {
	Config    *configs.AppConfig
	Telemetry observability.Telemetry
}

// PingDatabase checks the database health by establishing a connection and pinging it.
func (c *HealthConsumerImpl) PingDatabase(ctx context.Context) error {
	// Log when the request came in
	c.Telemetry.LogInfo("Received PingDatabase request", map[string]interface{}{
		"target_db": c.Config.DatabaseHost,
		"db_name":   c.Config.DatabaseName,
	})

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Config.DatabaseHost,
		c.Config.DatabasePort,
		c.Config.DatabaseUser,
		c.Config.DatabasePassword,
		c.Config.DatabaseName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		c.Telemetry.LogError("Failed to open database connection", err)
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		c.Telemetry.LogError("Failed to ping database", err)
		return err
	}

	// Log the response
	c.Telemetry.LogInfo("PingDatabase response", map[string]interface{}{
		"status":  200,
		"message": "Database is reachable and responding to ping",
	})

	return nil
}
