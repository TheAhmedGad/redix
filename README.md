# Redix

Redix is a Redis-like in-memory data store implementation in Go, providing a high-performance, lightweight alternative to Redis with a focus on simplicity and extensibility.

## Features

- RESP (Redis Serialization Protocol) implementation
- In-memory data storage with support for basic data types
- Pub/Sub messaging system
- Authentication support
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
│   ├── auth/              # Authentication implementation
│   ├── client/            # Client connection handling
│   ├── protocol/          # RESP protocol implementation
│   ├── pubsub/            # Pub/Sub messaging system
│   └── server/            # Server implementation
├── test/                  # Test suite
├── dockit/               # Docker-related files
│   ├── GoLang/          # Go service Dockerfile
│   └── mysql/           # MySQL initialization scripts
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

### Basic Commands

The server implements a subset of Redis commands:

- `SET key value` - Set a key-value pair
- `GET key` - Get the value of a key
- `DEL key` - Delete a key
- `PUBLISH channel message` - Publish a message to a channel
- `SUBSCRIBE channel` - Subscribe to a channel

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

- Inspired by Redis
- Built with Go
- Uses RESP protocol specification 