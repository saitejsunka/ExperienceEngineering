package observability

// Telemetry provides methods to record metrics and logs.
type Telemetry interface {
	LogInfo(message string, payload map[string]interface{})
	LogError(message string, err error)
	Close() error
}
