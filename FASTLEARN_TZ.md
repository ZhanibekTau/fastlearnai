# FastLearn — Technical Specification
## AI Summary для видеоуроков (EN/RU/KK)

---

## 1. OVERVIEW

### 1.1 Цель проекта
Создать сервис, который превращает видеоуроки, вебинары и подкасты в структурированные конспекты за 5-10 минут. Мультиязычный (EN/RU/KK), дешёвый ($2-3/мес), быстрый и специализированный под образование.

### 1.2 MVP Scope
**Что входит:**
- Web приложение (Vue.js)
- Backend API (Go)
- Базовая обработка видео
- Фримиум модель
- 3 языка (EN/RU/KK)

**Что НЕ входит (Phase 2):**
- Flutter мобилка (добавим после)
- YouTube Playlist автоматизация
- LMS интеграции
- Advanced Q&A генерация

### 1.3 Целевые пользователи (MVP)
1. YouTube-каналы с образовательным контентом
2. Студенты, которые смотрят онлайн-курсы
3. Преподаватели онлайн-школ
4. Контент-мейкеры (нужны саммери для видео)

---

## 2. FUNCTIONAL REQUIREMENTS

### 2.1 User Flows

#### Flow 1: Регистрация и вход
```
1. Пользователь переходит на fastlearn.ai
2. Может загрузить видео БЕЗ регистрации (1 бесплатный)
3. Или создать аккаунт (Google/GitHub + Email)
4. После регистрации:
   - 3 видео/мес бесплатно (free tier)
   - Платные тарифы доступны
```

#### Flow 2: Загрузка видео и генерация саммери
```
1. Пользователь нажимает "Upload Video"
2. Выбирает источник:
   - Загрузить файл (MP4, MOV, WebM, max 500MB)
   - YouTube ссылка
   - Вставить аудио файл (MP3, M4A, max 200MB)
3. Указывает язык видео (EN/RU/KK или авто-определение)
4. Система обрабатывает:
   - Whisper: транскрипция на том же языке
   - Claude 3.5 Sonnet: генерирует саммери
   - Парсинг структуры (заголовки, key takeaways, Q&A)
5. Результат выдаётся в 3 форматах:
   - HTML (в приложении)
   - PDF (скачать)
   - Markdown (для блога)
6. История сохраняется в аккаунте
```

#### Flow 3: Просмотр истории и повторная обработка
```
1. Пользователь видит список загруженных видео
2. Может:
   - Пересмотреть саммери
   - Скачать PDF/Markdown
   - Удалить видео
   - Поделиться ссылкой (публичная версия)
```

#### Flow 4: Подписка и платежи
```
1. Пользователь исчерпал лимит (free tier)
2. Видит "Upgrade to Pro"
3. Выбирает тариф:
   - Free: 3 видео/мес, до 30 мин каждое, только текст
   - Pro: $2.99/мес, 30 видео/мес, до 3 часов каждое, PDF + Markdown
   - Premium: $7.99/мес, неограниченно, API доступ
4. Оплата Stripe (карта)
5. Подписка активируется сразу
```

### 2.2 Функции приложения

| Функция | Free | Pro | Premium |
|---------|------|-----|---------|
| Видео/мес | 3 | 30 | Неограниченно |
| Макс длина видео | 30 мин | 3 часа | Неограниченно |
| Формат вывода | HTML | HTML + PDF + MD | + API |
| Языки | EN/RU/KK | EN/RU/KK | EN/RU/KK |
| История | 7 дней | 30 дней | Вечно |
| Q&A генерация | ❌ | ✅ | ✅ |
| Экспорт в JSON | ❌ | ❌ | ✅ |

### 2.3 Структура саммери (output)

```
{
  "title": "Название видео (или первая строка транскрипта)",
  "duration": "45:30",
  "language": "en",
  "summary": {
    "overview": "Краткое описание видео в 2-3 предложении",
    "key_points": [
      {
        "topic": "What is Machine Learning?",
        "explanation": "Детальное объяснение (2-3 параграфа)"
      },
      ...
    ],
    "key_takeaways": [
      "Вывод 1",
      "Вывод 2",
      "Вывод 3"
    ],
    "examples": [
      {
        "topic": "Linear Regression",
        "example": "Пример из видео"
      }
    ]
  },
  "qa": {
    "questions": [
      {
        "question": "Вопрос для самопроверки",
        "answer": "Ответ из видео"
      }
    ]
  },
  "transcript": "Полная транскрипция (если Premium или Pro)"
}
```

---

## 3. АРХИТЕКТУРА

### 3.1 System Overview

```
┌─────────────────┐
│   Vue.js Web    │
└────────┬────────┘
         │ HTTPS
    ┌────▼────────────────────┐
    │   Go Backend API        │
    │  (3 instances in K8s)   │
    └────┬─────────┬──────────┘
         │         │
    ┌────▼──┐  ┌───▼─────────────────┐
    │PostgreSQL│ Redis (Task Queue)  │
    └─────────┘ └───┬────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
    ┌───▼──────┐ ┌─▼────────┐ ┌─▼──────────┐
    │ Whisper  │ │  Claude  │ │ GCS (Files)│
    │ (Docker) │ │   API    │ │            │
    └──────────┘ └──────────┘ └────────────┘
```

### 3.2 Backend (Go)

#### 3.2.1 API Endpoints

```
# Auth
POST   /api/v1/auth/register        - Регистрация (Email, Google, GitHub)
POST   /api/v1/auth/login           - Вход
POST   /api/v1/auth/refresh         - Обновить JWT token
POST   /api/v1/auth/logout          - Выход

# Video Management
POST   /api/v1/videos/upload        - Загрузить видео / ссылку
GET    /api/v1/videos               - Список видео пользователя (pagination)
GET    /api/v1/videos/:id           - Деталь видео + саммери
DELETE /api/v1/videos/:id           - Удалить видео
GET    /api/v1/videos/:id/download  - Скачать PDF/MD

# Processing
GET    /api/v1/videos/:id/status    - Статус обработки (polling или WebSocket)
POST   /api/v1/videos/:id/regenerate - Пересоздать саммери

# Subscription
GET    /api/v1/user/subscription    - Информация о подписке
POST   /api/v1/subscription/upgrade - Upgrade на Pro/Premium
GET    /api/v1/subscription/plans   - Доступные тарифы

# Billing
POST   /api/v1/stripe/webhook       - Stripe callback для платежей
GET    /api/v1/user/usage           - Остаток видео в этом месяце

# Public
GET    /api/v1/videos/:id/public    - Публичная ссылка на саммери
GET    /health                      - Health check
```

#### 3.2.2 Go Project Structure

```
fastlearn-backend/
├── main.go                    # Entry point
├── config/
│   ├── config.go             # Config loading (env vars)
│   └── constants.go          # Constants
├── api/
│   ├── routes.go             # Route definitions
│   ├── handlers/
│   │   ├── auth.go           # Auth endpoints
│   │   ├── video.go          # Video upload/list
│   │   ├── processing.go     # Processing status
│   │   └── subscription.go   # Payment endpoints
│   └── middleware/
│       ├── auth.go           # JWT verification
│       ├── ratelimit.go      # Rate limiting
│       └── cors.go           # CORS
├── services/
│   ├── auth_service.go       # Auth logic
│   ├── video_service.go      # Video business logic
│   ├── whisper_service.go    # Whisper integration
│   ├── claude_service.go     # Claude API integration
│   ├── stripe_service.go     # Stripe integration
│   └── storage_service.go    # GCS upload/download
├── models/
│   ├── user.go               # User struct
│   ├── video.go              # Video struct
│   ├── subscription.go       # Subscription struct
│   └── summary.go            # Summary output struct
├── db/
│   ├── migrations/           # SQL migrations (goose/migrate)
│   │   └── 001_init.sql      # Initial schema
│   └── queries/
│       ├── user.go           # User queries (sqlc)
│       ├── video.go          # Video queries
│       └── subscription.go   # Subscription queries
├── queue/
│   ├── job.go                # Job structure
│   └── processor.go          # Redis queue processor
├── utils/
│   ├── jwt.go                # JWT utilities
│   ├── errors.go             # Error handling
│   └── logger.go             # Logging (structured)
└── docker-compose.yml        # Local dev (Postgres, Redis)
```

#### 3.2.3 Key Technologies

- **Framework:** Echo (fast HTTP framework)
- **Database:** sqlc (type-safe SQL) + goose (migrations)
- **Auth:** JWT tokens + Google OAuth / GitHub OAuth
- **Queue:** Redis + background worker (Goroutines)
- **Logging:** Structured logging (slog или zap)
- **HTTP Client:** Standard library + retry logic
- **Email:** SendGrid (для notifications)

### 3.3 Frontend (Vue.js)

#### 3.3.1 Pages/Components

```
src/
├── main.js
├── App.vue
├── router/
│   └── index.js              # Routing config
├── views/
│   ├── HomePage.vue          # Лендинг / Upload видео
│   ├── LoginPage.vue         # Вход / Регистрация
│   ├── DashboardPage.vue     # История видео
│   ├── VideoDetailPage.vue   # Саммери конкретного видео
│   ├── PricingPage.vue       # Тарифы
│   └── AccountPage.vue       # Профиль и подписка
├── components/
│   ├── VideoUpload.vue       # Upload форма
│   ├── SummaryViewer.vue     # Просмотр саммери
│   ├── VideoList.vue         # Таблица видео
│   ├── Header.vue            # Navigation
│   └── LoadingSpinner.vue    # Spinner при обработке
├── stores/
│   ├── auth.js               # Auth store (Pinia)
│   ├── video.js              # Video store
│   └── subscription.js       # Subscription store
└── api/
    └── client.js             # Axios instance + interceptors
```

#### 3.3.2 Key Features (Frontend)

- **Upload:** Drag-and-drop + file input
- **Real-time status:** WebSocket или polling (GET /status)
- **Export:** Download PDF, Copy markdown
- **Responsive:** Mobile-friendly (Bootstrap Vue или Tailwind)
- **Dark mode:** Toggle
- **Language:** EN/RU/KK (i18n)

### 3.4 Database (PostgreSQL)

#### 3.4.1 Schema

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),  -- NULL if OAuth
    name VARCHAR(255),
    oauth_provider VARCHAR(50),  -- 'google' или 'github'
    oauth_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    plan VARCHAR(50),  -- 'free', 'pro', 'premium'
    stripe_subscription_id VARCHAR(255),
    status VARCHAR(50),  -- 'active', 'cancelled', 'past_due'
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Videos table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    original_filename VARCHAR(255),
    storage_key VARCHAR(500),  -- GCS path
    duration_seconds INT,
    language VARCHAR(10),  -- 'en', 'ru', 'kk'
    file_size_bytes BIGINT,
    status VARCHAR(50),  -- 'processing', 'completed', 'failed'
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    INDEX(user_id, created_at)
);

-- Summaries table
CREATE TABLE summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    content JSONB,  -- Full summary JSON structure
    transcript TEXT,  -- Full transcript (nullable for Free tier)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(video_id)
);

-- Monthly usage tracking
CREATE TABLE usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    month DATE,  -- First day of month
    videos_count INT DEFAULT 0,
    total_minutes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, month)
);
```

#### 3.4.2 Indexes
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_videos_user_id_created ON videos(user_id, created_at DESC);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_usage_user_month ON usage(user_id, month);
```

### 3.5 Processing Architecture

#### 3.5.1 Async Job Queue (Redis + Worker)

```
User Upload Video
       ↓
   Go API (HTTP POST)
       ↓
   Validate & Store in DB (status: processing)
       ↓
   Push job to Redis Queue
       ↓
   Return response to user immediately
       ↓
   Background Worker (Goroutines):
       1. Download video from GCS
       2. Call Whisper API (transcription)
       3. Call Claude API (summary generation)
       4. Parse & structure output
       5. Store in DB (status: completed)
       6. Send email notification
```

#### 3.5.2 Job Structure

```go
type Job struct {
    VideoID     string
    UserID      string
    Language    string
    StorageKey  string
    FileSize    int64
    Priority    int  // 1-10
}
```

### 3.6 External Integrations

#### 3.6.1 Whisper (Transcription)
- **Option A (Recommended):** OpenAI Whisper API ($0.006/min)
  - Simple, no maintenance
  - Supports EN/RU/KK
  
- **Option B:** Self-hosted Docker container
  - More control, cheaper at scale
  - Harder to maintain
  
**Decision:** Use Whisper API for MVP, migrate to self-hosted later

#### 3.6.2 Claude (Summary Generation)
- **API:** Claude 3.5 Sonnet
- **Cost:** ~$0.003 per summary (100K tokens input, 1K output)
- **Timeout:** 30 seconds max
- **Retry:** 3 attempts with exponential backoff

#### 3.6.3 Storage (GCS - Google Cloud Storage)
- Store video files temporarily (24 hours, then delete)
- Store PDFs permanently (for downloads)
- Cost: ~$0.02/GB/month

#### 3.6.4 Stripe (Payments)
- Monthly recurring billing
- Webhook for subscription events (upgraded, cancelled, failed)
- Supports international cards

#### 3.6.5 SendGrid (Email)
- Transactional emails (reset password, summary ready, etc.)
- Free tier: 100/day

### 3.7 Docker & Deployment

#### 3.7.1 Dockerfile (Go Backend)

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o fastlearn-api ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/fastlearn-api /usr/local/bin/
EXPOSE 8080
CMD ["fastlearn-api"]
```

#### 3.7.2 Docker Compose (Local Development)

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: fastlearn
      POSTGRES_PASSWORD: dev_password
      POSTGRES_DB: fastlearn_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://fastlearn:dev_password@postgres:5432/fastlearn_db?sslmode=disable
      REDIS_URL: redis://redis:6379
      PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

volumes:
  postgres_data:
```

#### 3.7.3 Production Deployment

- **Hosting:** Google Cloud Run (auto-scaling) или DigitalOcean (K8s)
- **Database:** Cloud SQL (Postgres) или self-managed
- **Cache:** Cloud Memorystore (Redis)
- **Storage:** Google Cloud Storage
- **DNS:** Cloudflare
- **CI/CD:** GitHub Actions → Cloud Run

---

## 4. PERFORMANCE & SCALING

### 4.1 Latency Targets

| Operation | Target | Details |
|-----------|--------|---------|
| Video upload | < 2s | Just file upload |
| Job queuing | < 1s | DB + Redis write |
| User receives response | < 2s | API response time |
| Whisper transcription | 5-15 min | Depends on video length |
| Claude summary | 10-30s | After transcription ready |
| PDF generation | 5s | From summary data |
| **Total end-to-end** | **20-30 min** | For 1-hour video |

### 4.2 Capacity Planning

**Per instance (1 Go API pod):**
- ~1000 concurrent connections
- ~100 simultaneous job processing (with 4 worker threads)
- ~2000 req/sec peak

**Infrastructure (MVP):**
- 1x Go API instance (2 CPU, 2GB RAM) - $20/month
- 1x Postgres (shared, 2 CPU, 4GB) - $15/month
- 1x Redis (shared, 2GB) - $10/month
- GCS bucket - $5/month (1GB storage)
- Total: ~$50/month fixed costs

### 4.3 Scaling Strategy

**Phase 1 (0-100 users):** Single instance
**Phase 2 (100-1000 users):** Load balancer + 2x backend instances
**Phase 3 (1000+ users):** K8s cluster, auto-scaling

---

## 5. SECURITY

### 5.1 Authentication & Authorization

- JWT tokens (HS256, 24-hour expiry)
- Refresh tokens (7-day expiry, stored in secure HttpOnly cookies)
- OAuth2 (Google, GitHub)
- Rate limiting: 100 req/min per IP, 1000 req/min per user

### 5.2 Data Protection

- HTTPS only (TLS 1.3)
- Database encryption at rest (Cloud SQL default)
- Video files deleted after 24 hours
- PDFs encrypted (password-protected if user requests)
- No logs of API responses with user data

### 5.3 Payment Security

- PCI compliance via Stripe (no card data in our DB)
- Webhook signature verification (Stripe secret key)
- Idempotency keys for payment requests

---

## 6. MONITORING & LOGGING

### 6.1 Metrics

```
- API latency (p50, p95, p99)
- Error rate (4xx, 5xx)
- Queue depth (Redis job count)
- Whisper API latency
- Claude API latency
- User signup rate
- Subscription conversion
- Video processing success rate
```

### 6.2 Logging

- Structured JSON logs (slog)
- Log levels: DEBUG, INFO, WARN, ERROR
- ELK stack or Google Cloud Logging
- Alert on: Error rate > 5%, API latency > 5s, Queue depth > 1000

### 6.3 Uptime SLA

Target: 99.5% uptime (allows ~4 hours downtime/month)

---

## 7. FUTURE (Phase 2+)

### 7.1 Flutter Mobile App

```
Screen structure:
├── AuthScreen (Login/Register)
├── HomeScreen (Upload + List)
├── VideoDetailScreen (Summary view)
├── SettingsScreen (Profile, subscription)
└── DownloadScreen (PDF/Markdown)

Features:
- Camera + file picker for video upload
- Offline caching of summaries
- Push notifications (video processing done)
- Dark mode
```

### 7.2 Advanced Features

- YouTube playlist auto-summarize
- LMS integrations (Moodle, Canvas, Blackboard)
- Slack/Discord bot (share summaries)
- Browser extension (highlight in transcript)
- Multi-language learning (transcript side-by-side)
- Custom summary templates (different formats)

---

## 8. DEVELOPMENT TIMELINE

### Week 1 (Backend MVP)
- [ ] Go API scaffolding (Echo framework)
- [ ] Auth (JWT, Google OAuth)
- [ ] Postgres schema + migrations
- [ ] Video upload to GCS
- [ ] Redis queue setup
- [ ] Health checks + logging

### Week 2 (Processing)
- [ ] Whisper API integration
- [ ] Claude API integration
- [ ] Job processor (async worker)
- [ ] Summary JSON structure
- [ ] Email notifications

### Week 3 (Frontend)
- [ ] Vue.js setup + routing
- [ ] Auth pages (login, register)
- [ ] Upload component
- [ ] Video list page
- [ ] Summary viewer
- [ ] API client + error handling

### Week 4 (Payments & Polish)
- [ ] Stripe integration
- [ ] Subscription logic
- [ ] Usage tracking
- [ ] PDF export
- [ ] Markdown export
- [ ] Testing + bug fixes
- [ ] Deployment to Cloud Run
- [ ] Domain + HTTPS

### Week 5+ (Flutter & Growth)
- [ ] Flutter project setup
- [ ] Auth integration
- [ ] Video upload
- [ ] Summary view
- [ ] App Store + Play Store
- [ ] Beta testing with real users

---

## 9. TECH STACK SUMMARY

| Layer | Technology | Why |
|-------|-----------|-----|
| Backend | Go + Echo | Fast, concurrent, easy deployment |
| Frontend | Vue.js + Axios | Fast prototyping, good docs |
| Mobile | Flutter | iOS + Android simultaneously |
| Database | PostgreSQL + sqlc | Type-safe, migrations, relationships |
| Cache | Redis | Job queue, session storage |
| Storage | Google Cloud Storage | Simple, reliable, cheap |
| Transcription | OpenAI Whisper API | Best quality, supports all languages |
| Summary | Claude 3.5 Sonnet | Best quality, low cost |
| Payments | Stripe | Industry standard |
| Hosting | Google Cloud Run | Auto-scaling, pay-per-use |
| CI/CD | GitHub Actions | Free, integrated with GitHub |

---

## 10. SUCCESS METRICS (MVP Phase)

- **10 signups in 2 weeks** → Idea validation
- **3-5 paying customers by week 4** → Business viability
- **Average processing time < 20 min** → Performance acceptable
- **<1% error rate** → Quality acceptable
- **>80% retention after 7 days** → Product-market fit signal

---

# FastLearn MVP - Payment Strategy Update

## Phase 1: MVP (Current - Telegram based)

### Why Telegram?

✅ **Fast to implement** - no integration needed
✅ **Works right now** - Kaspi/PayBox require verification
✅ **Validates demand** - see who actually wants to pay
✅ **No compliance issues** - manual process, safer for MVP
✅ **Easy to scale** - Phase 2 can add automated payments

---

## How It Works

### User Flow

1. **User hits limit** (3 videos/month on free tier)
   ```
   Error: "Monthly limit reached (3/3)"
   "Want more? Contact us: t.me/fastlearn_support"
   ```

2. **User clicks Telegram link**
   ```
   Opens: https://t.me/fastlearn_support
   Pre-filled: "Hi! I want to upgrade my FastLearn plan"
   ```

3. **User sends message**
   ```
   "Hi! I want Pro plan"
   ```

4. **You respond in Telegram**
   ```
   "Great! Pro plan is $2.99/month. 
    Send me your email and I'll activate it."
   ```

5. **You update database** (one SQL query)
   ```sql
   UPDATE subscriptions 
   SET plan='pro', status='active'
   WHERE user_id='their_user_id'
   ```

6. **User can upload 30 videos/month** ✅

---

## Implementation Details

### Backend Changes (Already Done ✅)

**File:** `backend/services/video.go`
```go
// MVP: All users have free tier limit of 3 videos/month
freeLimit := 3

if usage.VideosCount >= freeLimit {
    return false, "Monthly limit reached. Contact us via Telegram", nil
}
```

**File:** `backend/api/handler.go`
```go
// New endpoint for upgrade info
func (h *Handler) GetUpgradeInfo(c echo.Context) error {
    return c.JSON(200, map[string]interface{}{
        "contact": map[string]string{
            "type": "telegram",
            "username": "@fastlearn_support",  // ← CHANGE THIS
            "url": "https://t.me/fastlearn_support",  // ← CHANGE THIS
        },
        "plans": []map[string]string{
            {"name": "Pro", "price": "$2.99/month", "videos": "30"},
            {"name": "Premium", "price": "$7.99/month", "videos": "Unlimited"},
        },
    })
}
```

### Frontend Changes (Already Done ✅)

**File:** `frontend/src/views/HomePage.vue`
```vue
<div v-if="usage?.videos_count >= 3" class="upgrade-notice">
    <h3>🎯 Monthly limit reached</h3>
    <a :href="telegramLink" target="_blank" class="btn-upgrade-telegram">
        💬 Upgrade via Telegram
    </a>
</div>
```

---

## Configuration

### Step 1: Create Telegram Bot (optional)
If you want a bot to automate responses later:

1. Message @BotFather in Telegram
2. Create new bot: `/newbot`
3. Name: FastLearn Support
4. Username: `@fastlearn_support_bot`
5. Copy token: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`

### Step 2: Update Code

**In `backend/api/handler.go` and `frontend/src/views/HomePage.vue`:**

Change:
```
@fastlearn_support  → @YOUR_USERNAME
https://t.me/fastlearn_support → https://t.me/YOUR_USERNAME
```

### Step 3: Manual Process (MVP)

Keep notifications in Telegram:
1. Telegram personal account: @fastlearn_support
2. Join FastLearn Telegram group / channel for notifications
3. When user messages, respond and update DB

---

## Database Updates (Manual)

### Pro Plan Activation

```sql
UPDATE subscriptions 
SET plan='pro', 
    status='active',
    current_period_start=CURRENT_TIMESTAMP,
    current_period_end=CURRENT_TIMESTAMP + INTERVAL '30 days'
WHERE user_id='user_id_here';
```

### Premium Plan Activation

```sql
UPDATE subscriptions 
SET plan='premium', 
    status='active',
    current_period_start=CURRENT_TIMESTAMP,
    current_period_end=CURRENT_TIMESTAMP + INTERVAL '30 days'
WHERE user_id='user_id_here';
```

### Cancel Subscription

```sql
UPDATE subscriptions 
SET plan='free', status='active'
WHERE user_id='user_id_here';
```

---

## Pricing Tiers (MVP)

| Plan | Price | Videos/mo | Duration | Support |
|------|-------|-----------|----------|---------|
| **Free** | $0 | 3 | 30 min max | Telegram |
| **Pro** | $2.99/mo | 30 | 3 hours max | Telegram |
| **Premium** | $7.99/mo | Unlimited | Unlimited | Telegram |

---

## Limitations & Next Steps

### MVP Limitations
❌ No automatic billing
❌ No payment history in app
❌ No invoice generation
❌ Manual activation only

### Phase 2 (Automated Payments)

Once you have paying customers:
1. Integrate PayBox API
2. Add webhooks for payment callbacks
3. Automatic subscription activation
4. Invoice/receipt generation
5. Recurring billing

---

## Conversation Example

```
User: Hi! I want to upgrade my FastLearn plan

You: Great! Which plan would you like?
   - Pro: 30 videos/month for $2.99/month
   - Premium: Unlimited for $7.99/month

User: I want Pro

You: Perfect! Please send me your Kaspi/PayBox details 
(or payment screenshot once they're set up in Phase 2)

You: (Update DB with SQL)
UPDATE subscriptions SET plan='pro' WHERE user_id='...';

You: Done! Your Pro plan is now active. 
You can upload 30 videos this month! 🎉

User: Thanks!
```

---

## Tracking Payments (MVP)

Create a simple spreadsheet:

| Date | Username | Email | Plan | Price | Notes |
|------|----------|-------|------|-------|-------|
| 2024-06-29 | user123 | email@test.com | Pro | $2.99 | Paid via Telegram |
| 2024-06-30 | user456 | email2@test.com | Premium | $7.99 | Paid via Telegram |

---

## Why This Works for MVP

1. **Validates demand** - You see who actually pays
2. **No risk** - You don't handle payments directly
3. **Compliant** - No payment processing regulations
4. **Fast** - No integration with Kaspi/PayBox
5. **Flexible** - Easy to change pricing/plans
6. **Personal** - Builds relationship with early users
7. **Safe** - You can refund if needed

---

## Transition to Phase 2

When you're ready for automated payments:

1. Kaspi / PayBox API integration
2. Webhook handling for payment confirmations
3. Automatic subscription activation
4. Billing history in app
5. Recurring billing setup

But for now, **keep it simple and manual**. You'll learn what users actually want.

---

**Remember:** The best feature is one that's shipped and used. This MVP approach gets you live faster! 🚀
