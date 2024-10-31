package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/zombocoder/go-redis-proxy/pkg/resp"
)

const CLIENT_BUFFER_SIZE = 4096             // Buffer size for client commands
const REDIS_BUFFER_SIZE = 32 * 1024 * 1024  // 32 MB buffer size for Redis responses
const CONNECTION_TIMEOUT = 30 * time.Second // Timeout for client connections

// Adjust buffer dynamically based on expected size
func AdjustBufferSize(size int) []byte {
	if size > REDIS_BUFFER_SIZE {
		log.Printf("Requested size %d exceeds max allowed by Redis\n", size)
		return make([]byte, REDIS_BUFFER_SIZE)
	}
	return make([]byte, size)
}

// Configuration structures for JSON parsing
type ServerConfig struct {
	Listen int           `json:"listen"`
	Master RedisServer   `json:"master"`
	Slave  []RedisServer `json:"slave"`
}

type RedisServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Define read-only Redis commands
var ReadOnlyCommands = []string{
	"info", "smembers", "hlen", "hmget", "srandmember", "hvals", "randomkey", "strlen",
	"dbsize", "keys", "ttl", "lindex", "type", "llen", "dump", "scard", "echo", "lrange",
	"zcount", "exists", "sdiff", "zrange", "mget", "zrank", "get", "getbit", "getrange",
	"zrevrange", "zrevrangebyscore", "hexists", "object", "sinter", "zrevrank", "hget",
	"zscore", "hgetall", "sismember",
}

// Check if a command is read-only
func IsReadOnlyCommand(command string) bool {
	for _, cmd := range ReadOnlyCommands {
		if strings.EqualFold(command, cmd) {
			return true
		}
	}
	return false
}

// Handle client connections
func HandleConnection(conn net.Conn, config ServerConfig) {
	defer conn.Close()

	masterConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Master.Host, config.Master.Port))
	if err != nil {
		log.Println("Error connecting to master:", err)
		return
	}
	defer masterConn.Close()

	slaveConns := make([]net.Conn, len(config.Slave))
	for i, slave := range config.Slave {
		slaveConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", slave.Host, slave.Port))
		if err != nil {
			log.Println("Error connecting to slave:", err)
			continue
		}
		slaveConns[i] = slaveConn
		defer slaveConn.Close()
	}

	clientBuf := make([]byte, CLIENT_BUFFER_SIZE)
	redisBuf := AdjustBufferSize(REDIS_BUFFER_SIZE)
	err = conn.SetReadDeadline(time.Now().Add(CONNECTION_TIMEOUT))
	if err != nil {
		log.Println("Error setting read deadline:", err)
		return
	}
	for {
		n, err := conn.Read(clientBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed the connection")
				break
			}
			log.Println("Error reading from client:", err)
			break
		}

		command, err := resp.ParseRESPCommand(clientBuf[:n])
		if err != nil {
			log.Println("Failed to parse command:", err)
			continue
		}

		log.Printf("Parsed command: '%s' from client '%s'\n", command, conn.RemoteAddr())

		var targetConn net.Conn
		if IsReadOnlyCommand(command) && len(slaveConns) > 0 {
			targetConn = slaveConns[0]
			log.Printf("Routing read-only command to slave: '%s' to '%s'\n", command, targetConn.RemoteAddr())
		} else {
			targetConn = masterConn
			log.Printf("Routing write command to master: '%s' to '%s'\n", command, targetConn.RemoteAddr())
		}

		err = targetConn.SetWriteDeadline(time.Now().Add(CONNECTION_TIMEOUT))
		if err != nil {
			log.Println("Error setting write deadline:", err)
			break
		}
		_, err = targetConn.Write(clientBuf[:n])
		if err != nil {
			log.Printf("Error writing to %s: %v\n", targetConn.RemoteAddr(), err)
			break
		}

		m, err := targetConn.Read(redisBuf)
		if err != nil {
			log.Printf("Error reading from %s: %v\n", targetConn.RemoteAddr(), err)
			break
		}
		response := string(redisBuf[:m])
		log.Printf("Received response from %s: %s\n", targetConn.RemoteAddr(), response)

		_, err = conn.Write(redisBuf[:m])
		if err != nil {
			log.Printf("Error sending response to client %s: %v\n", conn.RemoteAddr(), err)
			break
		}
		log.Printf("Sent response to client '%s': %s\n", conn.RemoteAddr(), response)
	}
}
