# Деплой fastlearn на сервер (вместо bloggers)

Этот документ — пошаговая инструкция для миграции сервера, на котором сейчас крутится
**bloggers** (systemd-сервисы `mediapro` + `threads-bot`, nginx отдаёт статику фронта), на
**fastlearn** (Go backend + Python FastAPI + Vue frontend + Postgres + Redis, всё в Docker).

Домен остаётся тот же — везде ниже замени `YOUR_DOMAIN` на реальный.

Все команды в разделах 1 и 2 выполняются вручную по SSH на сервере. Дальше (раздел 4)
деплой уже автоматический — через GitHub Actions.

---

## 0. Бэкап перед удалением (обязательно)

Прежде чем что-либо удалять — сними бэкап на всякий случай.

```bash
# Бэкап файлов bloggers
sudo tar -czvf /root/mediapro-backup-$(date +%Y%m%d).tar.gz /var/www/mediapro

# Бэкап базы bloggers (уточни реальное имя БД/юзера, см. blogger-social/.env на сервере)
sudo -u postgres pg_dump mediapro_db | gzip > /root/mediapro_db-backup-$(date +%Y%m%d).sql.gz
```

Скачай оба архива себе на локальную машину (`scp`), если сервер потом будет переустанавливаться.

---

## 1. Ручная очистка сервера от bloggers

### 1.1 Остановить и отключить сервисы

```bash
sudo systemctl stop mediapro threads-bot
sudo systemctl disable mediapro threads-bot
```

### 1.2 Удалить unit-файлы systemd

```bash
sudo rm -f /etc/systemd/system/mediapro.service
sudo rm -f /etc/systemd/system/threads-bot.service
sudo systemctl daemon-reload
```

### 1.3 Удалить nginx-конфиг bloggers

```bash
sudo rm -f /etc/nginx/sites-enabled/mediapro
sudo rm -f /etc/nginx/sites-available/mediapro
# nginx -t проверим после того, как положим конфиг fastlearn (шаг 2.4) —
# не перезагружай nginx, пока не появится новый server-блок для YOUR_DOMAIN,
# иначе домен временно останется без ответа.
```

### 1.4 Удалить файлы приложения

```bash
sudo rm -rf /var/www/mediapro
```

### 1.5 База данных bloggers (опционально)

Если `mediapro_db` больше не нужна (бэкап из шага 0 уже снят):

```bash
sudo -u postgres psql -c "DROP DATABASE mediapro_db;"
sudo -u postgres psql -c "DROP USER mediapro;"
```

Сам системный PostgreSQL (установленный через `apt` для bloggers) можно не трогать —
fastlearn использует свой Postgres в Docker-контейнере на отдельном порту внутри
docker-сети, конфликтовать с системным они не будут (см. шаг 2, порт наружу не
публикуется).

Если системный Postgres/Redis больше вообще не нужны:

```bash
sudo systemctl stop postgresql
sudo systemctl disable postgresql
```

---

## 2. Разовая настройка сервера под fastlearn

### 2.1 Установить Docker + Docker Compose plugin

```bash
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
# перелогинься по SSH, чтобы группа docker применилась
```

Проверка: `docker compose version` должна отработать без sudo.

### 2.2 Каталог деплоя и .env

```bash
sudo mkdir -p /opt/fastlearn
sudo chown $USER:$USER /opt/fastlearn
```

Создай `/opt/fastlearn/.env` с реальными продовыми значениями (на основе
`.env.example` из репозитория — `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`,
`OPENAI_API_KEY`, `GCS_BUCKET`, `FRONTEND_URL=https://YOUR_DOMAIN` и т.д.).
Этот файл **живёт только на сервере** и никогда не коммитится (он в `.gitignore`).

```bash
chmod 600 /opt/fastlearn/.env
```

`docker-compose.prod.yml` в workflow копируется в `/opt/fastlearn/` при каждом деплое,
а `.env` в этой же папке остаётся нетронутым — `docker compose` подхватит его
автоматически (Compose сам читает `.env` из рабочей директории).

### 2.3 SSH-ключ для GitHub Actions

Сгенерируй отдельную ключевую пару для деплоя (не обязательно, но чище, чем переиспользовать
ключ bloggers):

```bash
ssh-keygen -t ed25519 -f ~/fastlearn_deploy_key -N ""
cat ~/fastlearn_deploy_key.pub >> ~/.ssh/authorized_keys
```

Приватный ключ (`~/fastlearn_deploy_key`, без `.pub`) пойдёт в секрет `SSH_PRIVATE_KEY`
на GitHub (см. раздел 3). Сервер уже использует нестандартный SSH-порт `2222`
(см. `settings.md` bloggers) — в workflow он уже прописан.

### 2.4 Nginx + TLS

Скопируй `deploy/nginx-fastlearn.conf` из репозитория на сервер, замени `YOUR_DOMAIN`
на реальный домен:

```bash
sudo cp deploy/nginx-fastlearn.conf /etc/nginx/sites-available/fastlearn
sudo sed -i 's/YOUR_DOMAIN/реальный-домен/g' /etc/nginx/sites-available/fastlearn
sudo ln -s /etc/nginx/sites-available/fastlearn /etc/nginx/sites-enabled/fastlearn
sudo nginx -t && sudo systemctl reload nginx
```

Если сертификат Let's Encrypt для этого домена уже есть (остался от bloggers) —
пути `/etc/letsencrypt/live/YOUR_DOMAIN/...` те же, ничего переиздавать не нужно.
Если сертификата нет или он был для другого имени — `sudo certbot --nginx -d YOUR_DOMAIN`.

### 2.5 Первый запуск вручную (до первого GH Actions деплоя)

Собери и подними стек прямо с сервера один раз, чтобы убедиться, что всё поднимается
(дальше уже будет через workflow готовыми образами из GHCR):

```bash
cd /opt/fastlearn
# временно нужен docker-compose.prod.yml + .env — скопируй docker-compose.prod.yml сюда вручную
docker compose -f docker-compose.prod.yml pull   # упадёт, если образов ещё нет в GHCR — это нормально до первого push
```

Реальный первый успешный `pull` произойдёт после первого пуша в `main` (раздел 4).

---

## 3. Настройка репозитория GitHub

Секреты добавляются в `Settings → Secrets and variables → Actions` репозитория
`ZhanibekTau/fastlearn`:

| Secret | Значение |
|---|---|
| `SERVER_HOST` | IP или хост сервера |
| `SERVER_USER` | пользователь для SSH-деплоя |
| `SSH_PRIVATE_KEY` | приватный ключ из шага 2.3 |
| `APP_URL` | `https://YOUR_DOMAIN` (для health-check) |

Плюс создать GitHub Environment `production` (`Settings → Environments`) — на него
ссылается job `deploy` в workflow; можно оставить без доп. правил, либо повесить
required reviewers, если хочешь ручное подтверждение перед каждым деплоем.

### Видимость образов в GHCR

По умолчанию пакеты в GHCR для приватного репо создаются приватными — сервер не
сможет их `pull` без авторизации. Проще всего после первого пуска:

`GitHub → твой профиль → Packages → fastlearn-backend / fastlearn-fastapi /
fastlearn-frontend → Package settings → Change visibility → Public`.

(Это публикует только Docker-образы, не исходный код репозитория.) Если публиковать
образы публично не хочется — дай знать, добавлю в workflow шаг с `docker login` на
сервере по PAT-токену с правом `read:packages`.

---

## 4. Как это работает дальше

При каждом `git push` в `main` (или ручной запуск через `workflow_dispatch`):

1. Job `build` собирает 3 образа (backend, fastapi, frontend) и пушит их в
   `ghcr.io/zhanibektau/fastlearn-*` с тегами `latest` и `<git-sha>`.
2. Job `deploy` копирует свежий `docker-compose.prod.yml` на сервер в
   `/opt/fastlearn/`, затем по SSH делает `docker compose pull && up -d` с новыми
   образами и подчищает старые (`docker image prune -f`).
3. Health-check дергает `${APP_URL}/health`.

`.env` с секретами при этом не трогается — он один раз создан вручную на сервере
(шаг 2.2) и остаётся там постоянно.