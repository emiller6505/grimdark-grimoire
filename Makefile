.PHONY: install start start-api start-web dev test clean help

# Default target
help:
	@echo "Grimoire Development Commands:"
	@echo ""
	@echo "  make install    - Install all dependencies (Go + npm)"
	@echo "  make start      - Start both API and web (same terminal)"
	@echo "  make dev        - Start both API and web (separate terminals)"
	@echo "  make start-api  - Start only the API server"
	@echo "  make start-web  - Start only the web frontend"
	@echo "  make test       - Run all tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make update-submodule - Update wh40k-10e submodule to latest"

# Install all dependencies
install:
	@echo "Installing Go dependencies..."
	cd api && go mod download
	@echo "Installing npm dependencies..."
	cd web && npm install
	@echo "Initializing git submodule..."
	git submodule update --init --recursive
	@echo "✓ All dependencies installed"

# Start both services in same terminal
start:
	@./start.sh

# Start both services in separate terminals
dev:
	@./start-dev.sh

# Start only API server
start-api:
	@echo "Starting API server on http://localhost:8080"
	cd api && DATA_DIR=../wh40k-10e go run ./cmd/server

# Start only web frontend
start-web:
	@echo "Starting web frontend on http://localhost:3000"
	cd web && npm run dev

# Run all tests
test:
	@echo "Running Go tests..."
	cd api && go test ./...
	@echo "✓ All tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	cd api && rm -rf bin/ server
	cd web && rm -rf dist/ node_modules/.vite
	@echo "✓ Clean complete"

# Update submodule
update-submodule:
	@echo "Updating wh40k-10e submodule..."
	git submodule update --remote wh40k-10e
	@echo "✓ Submodule updated"
	@echo "Don't forget to commit: git add wh40k-10e && git commit -m 'Update submodule'"

