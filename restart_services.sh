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

# 1.5 Clean up data
echo -e "\n${GREEN}Step 1.5: Cleaning up data and logs...${NC}"
rm -f *.log
rm -f anomalies.log
echo "Logs cleaned."

# Clean up Redis (if running in docker)
if docker ps | grep -q pp-redis; then
    echo "Cleaning Redis..."
    docker exec pp-redis redis-cli FLUSHALL > /dev/null 2>&1
    echo "Redis cache cleared."
else
    echo "Redis container (pp-redis) not found, skipping."
fi

# Clean up Postgres (if running in docker)
if docker ps | grep -q pp-postgres; then
    echo "Cleaning Database..."
    # Truncate all tables
    docker exec pp-postgres psql -U postgres -d pp_db -c "TRUNCATE TABLE users, wallets, orders, wallet_ledgers, settlement_logs RESTART IDENTITY CASCADE;" > /dev/null 2>&1
    echo "Database tables truncated."
else
    echo "Postgres container (pp-postgres) not found, skipping."
fi

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
