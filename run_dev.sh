#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

# Инструменты
GO=/usr/local/go/bin/go
NODE_BIN=/home/user/.config/JetBrains/GoLand2025.1/node/versions/20.15.0/bin
export PATH=$PATH:/usr/local/go/bin:$NODE_BIN:/home/user/.local/bin

# Загружаем .env
set -a; source .env; set +a

echo "==> Starting postgres + redis..."
docker-compose up -d postgres redis

echo "==> Waiting for postgres..."
until docker-compose exec -T postgres pg_isready -U "$DB_USER" -q 2>/dev/null; do
  sleep 1
done
echo "    postgres ready ✓"

echo "==> Creating uploads dir..."
mkdir -p "$ROOT/uploads/video"

# FastAPI
echo ""
echo "==> Installing FastAPI deps..."
python3 -m pip install -q --user --break-system-packages -r "$ROOT/fastapi/requirements.txt"

echo "==> Starting FastAPI (port 8001)..."
cd "$ROOT/fastapi"
python3 -m uvicorn main:app --host 0.0.0.0 --port 8001 --reload &
FASTAPI_PID=$!
cd "$ROOT"

# Backend
echo "==> Starting Go backend (port 8080)..."
cd "$ROOT/backend"
$GO run ./cmd/main.go &
BACKEND_PID=$!
cd "$ROOT"

# Frontend
echo "==> Installing frontend deps..."
cd "$ROOT/frontend"
npm install -q
echo "==> Starting Vue frontend (port 5173)..."
npm run dev &
FRONTEND_PID=$!
cd "$ROOT"

echo ""
echo "╔══════════════════════════════════════╗"
echo "║        FastLearn is running!         ║"
echo "╠══════════════════════════════════════╣"
echo "║  Frontend  → http://localhost:5173   ║"
echo "║  Backend   → http://localhost:8080   ║"
echo "║  FastAPI   → http://localhost:8001   ║"
echo "║  Health    → http://localhost:8080/health ║"
echo "╚══════════════════════════════════════╝"
echo "Press Ctrl+C to stop all services"

trap "echo ''; echo 'Stopping...'; kill $FASTAPI_PID $BACKEND_PID $FRONTEND_PID 2>/dev/null; docker-compose stop postgres redis; exit 0" SIGINT SIGTERM

wait
