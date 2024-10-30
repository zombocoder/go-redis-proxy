package main

import (
	"log"
	"os"
	"time"

	"github.com/zombocoder/go-redis-proxy/internal/app/redis_proxy"
	"github.com/zombocoder/go-redis-proxy/internal/pkg/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config_file>")
	}
	// Read the configuration file
	configFile := os.Args[1]
	// Load the server configuration
	configs, err := config.ParseConfigFile(&configFile)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
	// Start the server
	redis_proxy.StartAllServers(configs)
	// Start the memory logger
	redis_proxy.StartMemoryLogger(10 * time.Second)

	// Prevent main from exiting
	select {}
}
