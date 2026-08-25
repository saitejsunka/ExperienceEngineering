package observability

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

// InitializeObservabilityAndTelemetry initializes the GCP Cloud Logging client.
// If it fails (e.g., due to lack of credentials locally), it falls back to a standard logger.
func InitializeObservabilityAndTelemetry(ctx context.Context, projectID string) (Telemetry, error) {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		// Log to standard error if we cannot initialize GCP logging
		log.Printf("Failed to create GCP logging client: %v. Falling back to local logging.", err)
		return &localTelemetry{}, nil
	}

	// Create a global logger stream for the entire application backend
	logger := client.Logger("dbexp-backend")

	return &gcpTelemetry{
		client: client,
		logger: logger,
	}, nil
}
