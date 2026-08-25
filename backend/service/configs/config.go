package configs

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// AppConfig holds all the configuration for the application
type AppConfig struct {
	GCPProjectID     string `json:"GCPProjectID"`
	DatabaseHost     string `json:"DatabaseHost"`
	DatabasePort     string `json:"DatabasePort"`
	DatabaseUser     string `json:"DatabaseUser"`
	DatabasePassword string `json:"DatabasePassword"`
	DatabaseName     string `json:"DatabaseName"`
}

// LoadConfigFromSecretManager fetches the configuration JSON from GCP Secret Manager
func LoadConfigFromSecretManager(ctx context.Context, projectID, secretName string) (*AppConfig, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	// Build the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName),
	}

	// Call the API
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	// Unmarshal the payload
	var config AppConfig
	if err := json.Unmarshal(result.Payload.Data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret payload: %w", err)
	}

	// Assign the project ID into the config as well
	config.GCPProjectID = projectID

	return &config, nil
}
