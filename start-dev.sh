#!/bin/bash

# Grimoire Development Startup Script (Separate Terminals)
# Opens separate terminal windows for API and web frontend

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

# Detect OS and open terminals accordingly
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS - use osascript to open new Terminal windows
    echo -e "${BLUE}Opening API server in new terminal...${NC}"
    osascript -e "tell application \"Terminal\" to do script \"cd '$(pwd)/api' && export DATA_DIR='../wh40k-10e' && go run ./cmd/server\""
    
    sleep 1
    
    echo -e "${BLUE}Opening web frontend in new terminal...${NC}"
    osascript -e "tell application \"Terminal\" to do script \"cd '$(pwd)/web' && npm run dev\""
    
    echo -e "${GREEN}✓ Services starting in separate terminals${NC}"
    echo -e "${YELLOW}API: http://localhost:8080${NC}"
    echo -e "${YELLOW}Web: http://localhost:3000${NC}"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux - try to detect desktop environment
    if command -v gnome-terminal &> /dev/null; then
        gnome-terminal -- bash -c "cd $(pwd)/api && export DATA_DIR='../wh40k-10e' && go run ./cmd/server; exec bash"
        gnome-terminal -- bash -c "cd $(pwd)/web && npm run dev; exec bash"
    elif command -v xterm &> /dev/null; then
        xterm -e "cd $(pwd)/api && export DATA_DIR='../wh40k-10e' && go run ./cmd/server" &
        xterm -e "cd $(pwd)/web && npm run dev" &
    else
        echo -e "${YELLOW}Could not detect terminal. Please run manually:${NC}"
        echo "Terminal 1: cd api && go run ./cmd/server"
        echo "Terminal 2: cd web && npm run dev"
    fi
else
    echo -e "${YELLOW}Unsupported OS. Please run manually:${NC}"
    echo "Terminal 1: cd api && go run ./cmd/server"
    echo "Terminal 2: cd web && npm run dev"
fi

