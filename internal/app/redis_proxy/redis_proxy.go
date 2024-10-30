package redis_proxy

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/zombocoder/go-redis-proxy/internal/pkg/server"
	"github.com/zombocoder/go-redis-proxy/pkg/stats"
)

// Start server based on configurations
func StartServer(config server.ServerConfig) {
	address := fmt.Sprintf(":%d", config.Listen)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting server on %s: %v", address, err)
	}
	defer listener.Close()
	log.Printf("Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go server.HandleConnection(conn, config)
	}
}

// StartAllServers starts all servers based on configurations
func StartAllServers(configs []server.ServerConfig) {
	for _, config := range configs {
		go StartServer(config)
	}
}

// StartMemoryLogger logs memory usage periodically
func StartMemoryLogger(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			stats.LogMemoryUsage()
		}
	}()
}
