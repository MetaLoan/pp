#!/bin/bash

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Restarting PP Services ===${NC}"

# Function to kill process on a port
kill_port() {
    local port=$1
    local name=$2
    echo -e "Checking $name on port $port..."
    pids=$(lsof -ti:$port)
    if [ -n "$pids" ]; then
        echo -e "${RED}Killing $name (PIDs: $pids)...${NC}"
        echo "$pids" | xargs kill -9
        echo "$name stopped."
    else
        echo "$name is not running."
    fi
}

# 1. Stop all services
echo -e "\n${GREEN}Step 1: Stopping services...${NC}"
kill_port 8080 "Backend"
kill_port 3000 "Frontend"

# Wait a moment to ensure ports are freed
sleep 2

# 2. Restart services
echo -e "\n${GREEN}Step 2: Starting services...${NC}"

# Start Backend
echo "Starting Backend (go run cmd/api/main.go)..."
nohup go run cmd/api/main.go > backend.log 2>&1 &
BACKEND_PID=$!
echo -e "Backend started with PID ${GREEN}$BACKEND_PID${NC}. Logs: backend.log"

# Start Frontend
echo "Starting Frontend (npm run dev)..."
cd web
nohup npm run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..
echo -e "Frontend started with PID ${GREEN}$FRONTEND_PID${NC}. Logs: frontend.log"

echo -e "\n${GREEN}=== Restart Complete ===${NC}"
echo "Backend API: http://localhost:8080"
echo "Frontend UI: http://localhost:3000"
echo "Check backend.log and frontend.log for details."
