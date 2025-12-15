# Grimoire

Warhammer 40,000 10th Edition API and Web Frontend

A full-stack application for browsing and searching Warhammer 40K unit data from BattleScribe XML files.

## Quick Start

### Install Dependencies

```bash
make install
# or
./start.sh  # (will install and start)
```

### Start Development Servers

**Option 1: Same Terminal (Background Processes)**
```bash
make start
# or
./start.sh
```

**Option 2: Separate Terminals (Recommended for Development)**
```bash
make dev
# or
./start-dev.sh
```

**Option 3: Start Individually**
```bash
make start-api   # API only (http://localhost:8080)
make start-web   # Web only (http://localhost:3000)
```

## Project Structure

```
grimoire/
├── api/              # Go backend API
├── web/              # React frontend
├── wh40k-10e/        # Data files (git submodule)
├── start.sh          # Start both services (same terminal)
├── start-dev.sh      # Start both services (separate terminals)
└── Makefile          # Convenience commands
```

## Available Commands

- `make install` - Install all dependencies
- `make start` - Start both API and web (same terminal)
- `make dev` - Start both API and web (separate terminals)
- `make start-api` - Start only API server
- `make start-web` - Start only web frontend
- `make test` - Run all tests
- `make clean` - Clean build artifacts
- `make update-submodule` - Update wh40k-10e submodule

## Environment Variables

### API
- `DATA_DIR` - Path to data files (default: `../wh40k-10e`)
- `PORT` - Server port (default: `8080`)
- `GIN_MODE` - Gin mode: `debug` or `release` (default: `debug`)

### Web
- `VITE_API_URL` - API base URL (default: `http://localhost:8080/api/v1`)

## Submodule Maintenance

The `wh40k-10e/` directory is a git submodule. See [SUBMODULE_MAINTENANCE.md](./SUBMODULE_MAINTENANCE.md) for details.

Quick update:
```bash
make update-submodule
git add wh40k-10e
git commit -m "Update submodule"
```

## Deployment

See `render.yaml` for Render deployment configuration. Both services are configured to deploy automatically.

## Documentation

- [API Documentation](./api/README.md)
- [Submodule Maintenance](./SUBMODULE_MAINTENANCE.md)

