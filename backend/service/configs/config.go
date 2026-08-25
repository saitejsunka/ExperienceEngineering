package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppConfig holds all the configuration for the application
type AppConfig struct {
	GCPProjectID     string `json:"GCPProjectID"`
	DatabaseHost     string `json:"DatabaseHost"`
	DatabasePort     int    `json:"DatabasePort"`
	DatabaseUser     string `json:"DatabaseUser"`
	DatabasePassword string `json:"DatabasePassword"`
	DatabaseName     string `json:"DatabaseName"`
}

// LoadConfig reads the .config file from the given path and parses it into AppConfig
func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config JSON: %w", err)
	}

	return &config, nil
}
