# Redix

Redix is a lightweight Redis-like pub/sub messaging system implemented in Go, providing a simple and efficient way to handle publish/subscribe patterns. It supports multiple authentication tokens, making it perfect for SaaS applications where different clients need isolated pub/sub channels.

## Features

- RESP (Redis Serialization Protocol) implementation
- Pub/Sub messaging system with token-based authentication
- Multi-tenant support through unique authentication tokens
- Channel isolation per token for secure SaaS deployments
- Simple and lightweight
- Docker support for easy deployment
- Comprehensive test suite

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for containerized deployment)
- Air (for live reload during development)

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/redix.git
cd redix
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o redix
```

### Using Docker

```bash
docker-compose up --build
```

## Development

The project uses Air for live reload during development. To start the development server:

```bash
air
```

## Project Structure

```
redix/
├── pkg/                    # Core packages
│   ├── client/            # Client connection handling
│   ├── protocol/          # RESP protocol implementation
│   ├── pubsub/            # Pub/Sub messaging system
│   └── server/            # Server implementation
├── test/                  # Test suite
├── dockit/               # Docker-related files
│   └── GoLang/          # Go service Dockerfile
├── main.go              # Application entry point
├── docker-compose.yml   # Docker Compose configuration
├── .air.toml           # Air live reload configuration
└── go.mod              # Go module definition
```

## Usage

### Starting the Server

```bash
./redix
```

By default, the server listens on `localhost:6379`.

### Authentication

Redix uses token-based authentication to support multiple tenants. Each token provides isolated access to pub/sub channels:

```bash
# Connect with a specific token
redis-cli -a your-token-here

# Or use the AUTH command after connecting
redis-cli
> AUTH your-token-here
```

### Pub/Sub Commands

The server implements Redis-style pub/sub commands with token-based isolation:

- `AUTH token` - Authenticate with a specific token
- `PUBLISH channel message` - Publish a message to a channel (scoped to the authenticated token)
- `SUBSCRIBE channel` - Subscribe to a channel (scoped to the authenticated token)
- `UNSUBSCRIBE channel` - Unsubscribe from a channel

Example usage with redis-cli:

```bash
# Terminal 1 - Subscribe to a channel with token1
redis-cli -a token1
> SUBSCRIBE mychannel

# Terminal 2 - Subscribe to the same channel with token2
redis-cli -a token2
> SUBSCRIBE mychannel

# Terminal 3 - Publish a message with token1
redis-cli -a token1
> PUBLISH mychannel "Hello from token1!"

# Terminal 4 - Publish a message with token2
redis-cli -a token2
> PUBLISH mychannel "Hello from token2!"
```

In this example, subscribers with token1 will only receive messages published with token1, and subscribers with token2 will only receive messages published with token2, even though they're using the same channel name. This isolation makes Redix suitable for SaaS applications where you need to keep different clients' messages separate.

## Testing

Run the test suite:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Redis pub/sub functionality
- Built with Go
- Uses RESP protocol specification 