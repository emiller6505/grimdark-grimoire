.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t grimoire-api:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-v $(PWD)/../wh40k-10e:/app/wh40k-10e:ro \
		-e DATA_DIR=/app/wh40k-10e \
		grimoire-api:latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

