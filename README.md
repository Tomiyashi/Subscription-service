# Subscription Service API 🚀

RESTful API для управления подписками пользователей, реализованный на Go. Проект демонстрирует чистую архитектуру (Repository → Service → Handler), работу с PostgreSQL через `pgx`, валидацию бизнес-логики и контейнеризацию через Docker.

## 🛠 Стек технологий

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.21+ |
| Роутер | `go-chi/chi/v5` |
| База данных | PostgreSQL 15 |
| Драйвер БД | `jackc/pgx/v5` |
| Контейнеризация | Docker & Docker Compose |
| Миграции | Custom Go-скрипт + SQL |

# 📁 Структура проекта
``` bash
subscription-service/
├── cmd/
│ ├── api/main.go # Точка входа HTTP-сервера
│ └── migrate/main.go # Скрипт применения миграций
├── internal/
│ ├── handlers/handler.go # HTTP-обработчики (Chi + JSON)
│ ├── models/subscription.go # DTO и доменные модели
│ ├── repository/ # Слой доступа к данным (pgxpool)
│ └── service/ # Бизнес-логика и валидация
├── migrations/ # SQL-скрипты структуры БД
├── configs/ # Конфигурации (заглушки)
├── docker-compose.yml # Оркестрация PostgreSQL
├── .env # Переменные окружения
├── go.mod / go.sum
└── README.md
```

# ⚙️ Требования

- Go 1.21 или новее
- Docker & Docker Compose
- Postman / curl / Insomnia (для тестирования)


# 🚀 Быстрый старт

## 1. Клонирование

``` bash
- git clone <твой-репозиторий>
- cd subscription-service
```

## 2. Настройка .env
Создай файл .env в корне проекта:
``` bash
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions_db
PORT=8080
```
## 3. Запуск
### Запуск базы данных
docker compose up -d

### Применение миграций
go run cmd/migrate/main.go

### Запуск сервера
go run cmd/api/main.go

Сервер запустится на http://localhost:8080


# 📡 API

## POST /subscriptions
Создать подписку
``` bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "YouTube Premium",
    "price": 199,
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "start_date": "2026-05-01T00:00:00Z"
  }'
```
## GET /subscriptions/{id}
Получить подписку по ID
``` bash
curl http://localhost:8080/subscriptions/<uuid>
``` 
## GET /subscriptions?user_id={uuid}
Получить список подписок пользователя

``` bash
curl "http://localhost:8080/subscriptions?user_id=<uuid>"
```

## 🏗 Архитектура
Handlers — HTTP слой (Chi router)
Service — Бизнес-логика и валидация
Repository — Работа с PostgreSQL (pgx)
Models — Структуры данных

## 📦 Зависимости
``` bash
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/google/uuid
```