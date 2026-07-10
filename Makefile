include .env
export

# ─── Инфраструктура ───────────────────────────────────────────────────────────

infra:                        ## Запустить только Postgres + Redis
	docker-compose up -d postgres redis
	@echo "Waiting for postgres..."
	@until docker-compose exec -T postgres pg_isready -U $(DB_USER) -q; do sleep 1; done
	@echo "Postgres ready"

infra-down:                   ## Остановить Postgres + Redis
	docker-compose stop postgres redis

# ─── Сервисы (локально) ───────────────────────────────────────────────────────

backend:                      ## Запустить Go backend (hot: air, иначе go run)
	@mkdir -p uploads/video
	cd backend && go run ./cmd/main.go

fastapi:                      ## Запустить Python FastAPI
	cd fastapi && uvicorn main:app --host 0.0.0.0 --port 8001 --reload

frontend:                     ## Запустить Vue.js dev server
	cd frontend && npm run dev

# ─── Все три сервиса в отдельных терминалах ──────────────────────────────────
# Запусти каждый в своём окне:
#   make infra
#   make backend   (новый терминал)
#   make fastapi   (новый терминал)
#   make frontend  (новый терминал)

# ─── Docker (полный стек) ─────────────────────────────────────────────────────

up:                           ## Запустить весь стек в Docker
	docker-compose up --build

up-d:                         ## Запустить весь стек в Docker (фон)
	docker-compose up --build -d

down:                         ## Остановить и удалить контейнеры
	docker-compose down

logs:                         ## Логи всех сервисов
	docker-compose logs -f

logs-back:                    ## Логи только backend
	docker-compose logs -f backend

logs-ai:                      ## Логи только fastapi
	docker-compose logs -f fastapi

# ─── Первичная настройка ──────────────────────────────────────────────────────

setup:                        ## Скопировать .env.example → .env и создать папки
	@test -f .env || (cp .env.example .env && echo "Created .env — заполни OPENAI_API_KEY!")
	@mkdir -p uploads/video
	@echo "uploads/video dir ready"

deps-backend:                 ## Скачать Go зависимости
	cd backend && go mod tidy

deps-fastapi:                 ## Установить Python зависимости
	cd fastapi && pip install -r requirements.txt

deps-frontend:                ## Установить Node зависимости
	cd frontend && npm install

deps: deps-backend deps-fastapi deps-frontend  ## Установить все зависимости

# ─── Утилиты ──────────────────────────────────────────────────────────────────

psql:                         ## Открыть psql в контейнере
	docker-compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)

redis-cli:                    ## Открыть redis-cli в контейнере
	docker-compose exec redis redis-cli

.PHONY: infra infra-down backend fastapi frontend up up-d down logs logs-back logs-ai setup deps deps-backend deps-fastapi deps-frontend psql redis-cli
