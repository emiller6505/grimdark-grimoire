#!/bin/bash

# Grimoire Development Startup Script
# Installs dependencies and starts both API and web frontend

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Grimoire Development Startup ===${NC}\n"

# Check if we're in the right directory
if [ ! -d "api" ] || [ ! -d "web" ]; then
    echo -e "${YELLOW}Error: api/ or web/ directory not found${NC}"
    echo "Please run this script from the repository root"
    exit 1
fi

# Install Go dependencies
echo -e "${GREEN}Installing Go dependencies...${NC}"
cd api
if command -v go &> /dev/null; then
    go mod download
    echo -e "${GREEN}✓ Go dependencies installed${NC}\n"
else
    echo -e "${YELLOW}Warning: Go not found, skipping Go dependencies${NC}\n"
fi
cd ..

# Install npm dependencies
echo -e "${GREEN}Installing npm dependencies...${NC}"
cd web
if command -v npm &> /dev/null; then
    npm install
    echo -e "${GREEN}✓ npm dependencies installed${NC}\n"
else
    echo -e "${YELLOW}Warning: npm not found, skipping npm dependencies${NC}\n"
fi
cd ..

# Check if submodule is initialized
if [ ! -f "wh40k-10e/Warhammer 40,000.gst" ]; then
    echo -e "${YELLOW}Initializing git submodule...${NC}"
    git submodule update --init --recursive
    echo -e "${GREEN}✓ Submodule initialized${NC}\n"
fi

# Set default DATA_DIR if not set
export DATA_DIR="${DATA_DIR:-../wh40k-10e}"
export PORT="${PORT:-8080}"

echo -e "${BLUE}Starting services...${NC}\n"
echo -e "${YELLOW}API will run on: http://localhost:${PORT}${NC}"
echo -e "${YELLOW}Web will run on: http://localhost:3000${NC}\n"
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}\n"

# Function to cleanup background processes
cleanup() {
    echo -e "\n${YELLOW}Stopping services...${NC}"
    kill $API_PID $WEB_PID 2>/dev/null || true
    wait $API_PID $WEB_PID 2>/dev/null || true
    echo -e "${GREEN}Services stopped${NC}"
    exit 0
}

# Trap Ctrl+C
trap cleanup SIGINT SIGTERM

# Start API server in background
echo -e "${BLUE}Starting API server...${NC}"
cd api
go run ./cmd/server &
API_PID=$!
cd ..

# Wait a moment for API to start
sleep 2

# Start web dev server in background
echo -e "${BLUE}Starting web frontend...${NC}"
cd web
npm run dev &
WEB_PID=$!
cd ..

# Wait for both processes
echo -e "${GREEN}✓ Both services started${NC}\n"
echo -e "${GREEN}API PID: ${API_PID}${NC}"
echo -e "${GREEN}Web PID: ${WEB_PID}${NC}\n"

# Wait for both processes
wait $API_PID $WEB_PID

