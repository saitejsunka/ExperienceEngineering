package observability

import (
	"fmt"
	"log"

	"cloud.google.com/go/logging"
)

type gcpTelemetry struct {
	client *logging.Client
	logger *logging.Logger
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
