package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zombocoder/go-redis-proxy/internal/pkg/server"
)

// parseConfigFile reads and parses the configuration file
func ParseConfigFile(configFile *string) ([]server.ServerConfig, error) {
	data, err := os.ReadFile(*configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var configs []server.ServerConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return configs, nil
}
