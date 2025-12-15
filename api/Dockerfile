FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files from api directory
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy source code from api directory
COPY api/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy data directory from repository root
COPY wh40k-10e ./wh40k-10e

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DATA_DIR=./wh40k-10e
ENV GIN_MODE=release

# Run the server
CMD ["./server"]

