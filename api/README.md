# Grimoire API

A Go REST API server for serving Warhammer 40,000 10th Edition BattleScribe XML data.

## Features

- Parse BattleScribe XML files (.cat and .gst)
- Resolve entry links to library entries
- Transform XML data to JSON
- In-memory caching for performance
- RESTful API endpoints for units, catalogues, factions, and search

## Prerequisites

- Go 1.21 or later
- Access to the `wh40k-10e` directory with XML files

## Installation

```bash
# Clone the repository
cd api

# Install dependencies
go mod download

# Build the application
make build
```

## Configuration

The API can be configured using environment variables:

- `DATA_DIR`: Path to the directory containing XML files (default: `../wh40k-10e`)
- `PORT`: Server port (default: `8080`)
- `GIN_MODE`: Gin mode - `debug` or `release` (default: `debug`)

## Running

### Local Development

```bash
# Set data directory
export DATA_DIR=../wh40k-10e

# Run the server
make run
# or
go run ./cmd/server
```

### Docker

```bash
# Build the image
make docker-build

# Run the container
make docker-run
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Game System
- `GET /api/v1/game-system` - Get game system information

### Catalogues
- `GET /api/v1/catalogues` - List all catalogues
- `GET /api/v1/catalogues/:id` - Get catalogue details
- `GET /api/v1/catalogues/:id/units` - Get units in a catalogue

### Units
- `GET /api/v1/units` - List units (with filters: `faction`, `category`, `search`, `limit`, `offset`)
- `GET /api/v1/units/:id` - Get unit details
- `GET /api/v1/units/:id/weapons` - Get unit weapons

### Factions
- `GET /api/v1/factions` - List all factions
- `GET /api/v1/factions/:name/units` - Get units by faction

### Search
- `GET /api/v1/search?q={query}&limit={limit}` - Search units

## Example Requests

```bash
# Get all units
curl http://localhost:8080/api/v1/units

# Get a specific unit
curl http://localhost:8080/api/v1/units/828d-840a-9a67-9074

# Search for units
curl http://localhost:8080/api/v1/search?q=marine

# Get units by faction
curl http://localhost:8080/api/v1/factions/Imperium/units
```

## Project Structure

```
api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── models/         # Data models
│   ├── parser/         # XML parsing logic
│   ├── handlers/       # HTTP handlers
│   ├── service/        # Business logic
│   └── cache/          # Caching layer
├── pkg/response/       # Response helpers
└── go.mod              # Go module file
```

## Development

```bash
# Format code
make fmt

# Run tests
make test

# Clean build artifacts
make clean
```

## License

This project uses BattleScribe data files which are maintained by the BSData community.

