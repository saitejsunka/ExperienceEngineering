package observability

import "log"

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
