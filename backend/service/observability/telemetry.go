package observability

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/logging"
)

// Telemetry provides methods to record metrics and logs.
type Telemetry interface {
	LogInfo(message string, payload map[string]interface{})
	LogError(message string, err error)
	Close() error
}

type gcpTelemetry struct {
	client *logging.Client
	logger *logging.Logger
}

// InitTelemetry initializes the GCP Cloud Logging client.
// If it fails (e.g., due to lack of credentials locally), it falls back to a standard logger.
func InitTelemetry(ctx context.Context, projectID string) (Telemetry, error) {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		// Log to standard error if we cannot initialize GCP logging
		log.Printf("Failed to create GCP logging client: %v. Falling back to local logging.", err)
		return &localTelemetry{}, nil
	}

	// Create a logger specifically for the database operations
	logger := client.Logger("dbexp-backend")

	return &gcpTelemetry{
		client: client,
		logger: logger,
	}, nil
}

// LogInfo records an informational message with structured payload.
func (t *gcpTelemetry) LogInfo(message string, payload map[string]interface{}) {
	if payload == nil {
		payload = make(map[string]interface{})
	}
	payload["message"] = message

	t.logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  payload,
	})
	log.Printf("[INFO] %s | %v", message, payload) // also print locally for dev
}

// LogError records an error message.
func (t *gcpTelemetry) LogError(message string, err error) {
	payload := map[string]interface{}{
		"message": message,
		"error":   err.Error(),
	}

	t.logger.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  payload,
	})
	log.Printf("[ERROR] %s: %v", message, err) // also print locally for dev
}

// Close gracefully flushes and closes the GCP logger.
func (t *gcpTelemetry) Close() error {
	if err := t.logger.Flush(); err != nil {
		return fmt.Errorf("failed to flush logs: %w", err)
	}
	return t.client.Close()
}

// --- Local Telemetry Fallback (for local development without GCP credentials) ---

type localTelemetry struct{}

func (l *localTelemetry) LogInfo(message string, payload map[string]interface{}) {
	log.Printf("[LOCAL INFO] %s | %v", message, payload)
}

func (l *localTelemetry) LogError(message string, err error) {
	log.Printf("[LOCAL ERROR] %s: %v", message, err)
}

func (l *localTelemetry) Close() error {
	return nil
}
