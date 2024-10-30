# Redis Proxy with Master-Slave Support

This project is a TCP-based Redis proxy written in Go. It forwards Redis commands to a master or slave server based on whether the command is read-only, allowing load distribution across master and slave servers.

## Features

- Routes read-only commands to Redis slave nodes to reduce load on the master.
- Routes write commands to the Redis master node to ensure data consistency.
- Configurable via a JSON file that specifies master and slave configurations.

## Configuration

The proxy reads its configuration from a JSON file, which specifies the port to listen on, the master Redis server, and any slave servers. Here’s an example configuration:

```json
[
	{
		"listen": 6380,
		"master": { "host": "host01", "port": 6378 },
		"slave": [{ "host": "host02", "port": 6380 }]
	},
	{
		"listen": 6381,
		"master": { "host": "host01", "port": 6379 },
		"slave": [
			{ "host": "host02", "port": 6380 },
			{ "host": "host03", "port": 6380 }
		]
	}
]
```

### Explanation of Config Fields

- `listen`: The port on which the proxy listens for client connections.
- `master`: Specifies the master Redis server (host and port).
- `slave`: Specifies one or more slave Redis servers (host and port).

## Usage

1. Clone this repository:

   ```bash
   git clone https://github.com/go-redis-proxy/redis-proxy.git
   cd redis-proxy
   ```

2. Build the project:

   ```bash
   go build -o redis-proxy ./cmd/server/main.go
   ```

3. Run the proxy with a configuration file:

   ```bash
   ./redis-proxy config/example_config.json
   ```

The proxy will start and listen on the specified ports as configured in `config.json`.

## How It Works

The proxy differentiates between read-only and write commands:

- **Read-only commands**: Commands that do not modify data (like `GET`, `INFO`, `SMEMBERS`, etc.) are routed to a slave node.
- **Write commands**: Commands that modify data (like `SET`, `INCR`, etc.) are routed to the master node.

The proxy uses an efficient lookup to determine whether a command is read-only, ensuring fast routing decisions.

## Read-Only Commands

The following commands are treated as read-only by the proxy (partial list):

- `GET`, `MGET`, `KEYS`, `EXISTS`, `SCARD`, `LENGTH`, `INFO`, etc.

The full list of read-only commands can be found in the source code.

## Contributing

Feel free to open issues or submit pull requests if you have suggestions, improvements, or bug fixes. Contributions are welcome!

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
