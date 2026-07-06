# Apex Karting Booking Service

REST API для онлайн-бронирования заездов картинг-центра.

Проект разработан в рамках выполнения тестового задания **SURF Summer School 2026** и представляет собой MVP backend-сервиса, реализующего полный цикл бронирования: авторизацию пользователя, просмотр доступных слотов, создание и отмену бронирований, а также управление профилем.

---

# Содержание

- [Описание проекта](#описание-проекта)
- [Основные возможности](#основные-возможности)
- [Архитектура](#архитектура)
- [Стек технологий](#стек-технологий)
- [Структура проекта](#структура-проекта)
- [Запуск проекта](#запуск-проекта)
- [Конфигурация](#конфигурация)
- [Структура базы данных](#структура-базы-данных)
- [REST API](#rest-api)
- [Тестирование](#тестирование)
- [Документация проекта](#документация-проекта)
- [Архитектурные решения](#архитектурные-решения)
- [Возможные улучшения](#возможные-улучшения)

---

# Описание проекта

Сервис предназначен для автоматизации процесса бронирования картинг-заездов.

Основной пользовательский сценарий:

1. Пользователь проходит авторизацию по OTP.
2. Получает JWT Access Token.
3. Просматривает доступные заезды.
4. Создает бронирование.
5. Управляет собственными бронированиями.
6. Просматривает и редактирует профиль.

Проект реализует исключительно backend-часть приложения и предоставляет REST API.

---

# Основные возможности

## Авторизация

- Авторизация по одноразовому OTP-коду
- Регистрация нового пользователя
- Авторизация существующего пользователя
- JWT Authentication
- Logout

---

## Работа со слотами

- Просмотр списка доступных заездов
- Просмотр информации о конкретном заезде
- Фильтрация по конфигурации трассы
- Отображение свободных мест
- Отображение доступной экипировки

---

## Управление профилем

- Просмотр профиля
- Обновление имени
- Обновление телефона
- Удаление аккаунта

---

## Управление бронированиями

- Создание бронирования
- Просмотр списка бронирований
- Просмотр информации о бронировании
- Отмена бронирования
- Поздняя отмена
- Защита от повторных запросов (Idempotency-Key)
- Защита от овербукинга

---

# Архитектура

Проект реализован с использованием многослойной архитектуры.

```
HTTP
│
▼
Handlers
│
▼
Services
│
▼
Repositories
│
▼
PostgreSQL
```

### Handler

Отвечает за:

- обработку HTTP-запросов;
- валидацию входных данных;
- сериализацию JSON;
- формирование HTTP-ответов.

---

### Service

Содержит бизнес-логику приложения.

Именно сервисы реализуют:

- создание бронирований;
- расчёт стоимости;
- валидацию данных;
- отмену бронирований;
- авторизацию.

---

### Repository

Инкапсулирует работу с PostgreSQL.

Весь SQL-код сосредоточен только в данном слое.

---

# Стек технологий

| Технология | Назначение |
|------------|------------|
| Go 1.25 | Backend |
| PostgreSQL | База данных |
| pgx/v5 | Работа с PostgreSQL |
| Chi Router | HTTP Router |
| JWT | Авторизация |
| Docker | Контейнеризация |
| Docker Compose | Локальная инфраструктура |

---

# Структура проекта

```
.
├── analysis
│
├── backend
│   ├── cmd
│   │   └── api
│   │
│   ├── internal
│   │   ├── app
│   │   ├── auth
│   │   ├── bookings
│   │   ├── config
│   │   ├── db
│   │   ├── http
│   │   │   ├── middleware
│   │   │   └── response
│   │   ├── profile
│   │   └── slots
│   │
│   ├── migrations
│   ├── docker-compose.yml
│   └── go.mod
│
└── docs
    └── testing
```

---

# Запуск проекта

## Требования

- Go 1.25+
- Docker
- Docker Compose

---

## 1. Клонирование

```bash
git clone <repository>
cd backend
```

---

## 2. Запуск PostgreSQL

```bash
docker compose up -d
```

---

## 3. Применение миграций

```bash
docker exec -i apex_karting \
psql -U karting -d karting \
< migrations/000001_init.sql
```

---

## 4. Запуск приложения

```bash
go run ./cmd/api
```

По умолчанию сервис запускается на

```
http://localhost:8080
```

---

# Конфигурация

Используется файл `.env`.

Пример:

```env
APP_PORT=8080

DATABASE_URL=postgres://karting:karting@localhost:5432/karting?sslmode=disable

JWT_SECRET=super_secret_key
```

---

# Структура базы данных

Основные сущности:

- users
- otp_codes
- track_configs
- marshals
- slots
- bookings
- booking_seats

Диаграмма базы данных находится в разделе документации.

---

# REST API

## Авторизация

```
POST /api/v1/auth/otp/send

POST /api/v1/auth/otp/verify

POST /api/v1/auth/logout
```

---

## Профиль

```
GET /api/v1/profile

PATCH /api/v1/profile

DELETE /api/v1/profile
```

---

## Слоты

```
GET /api/v1/slots

GET /api/v1/slots/{id}
```

---

## Бронирования

```
GET /api/v1/bookings

GET /api/v1/bookings/{id}

POST /api/v1/bookings

POST /api/v1/bookings/{id}/cancel
```

---

# Тестирование

Запуск всех тестов

```bash
go test ./...
```

Запуск отдельного пакета

```bash
go test ./internal/auth

go test ./internal/profile

go test ./internal/slots
```

Покрытие включает:

- unit-тесты сервисного слоя;
- ручные тест-кейсы.

---

# Документация проекта

В репозитории подготовлен полный комплект проектной документации.

## Аналитика

- Domain Model
- Business Requirements
- Functional Requirements
- Non-Functional Requirements
- User Stories
- Use Cases

---

## Проектирование

- ER Diagram
- Sequence Diagram
- OpenAPI Specification
- Screen Registry
- Design Brief

---

## Тестирование

- Unit Tests
- Manual Test Cases
- Bug Reports

---

# Архитектурные решения

## Repository Pattern

Изоляция SQL-кода от бизнес-логики.

---

## Service Layer

Бизнес-логика вынесена в сервисы.

---

## Dependency Injection

Все зависимости создаются при инициализации приложения.

---

## JWT Authentication

Защищённые маршруты используют middleware авторизации.

---

## PostgreSQL Transactions

Создание и отмена бронирования выполняются внутри одной транзакции.

Это обеспечивает консистентность данных.

---

## Idempotency-Key

Повторный запрос на создание бронирования не приводит к созданию дубликатов.

---

## Unit Testing

Сервисный слой покрыт unit-тестами с использованием mock-репозиториев.

---

# Возможные улучшения

В рамках MVP не реализованы:

- Refresh Token
- RBAC
- Swagger UI
- Redis для хранения OTP
- Rate Limiting
- Structured Logging
- Prometheus Metrics
- CI/CD Pipeline
- Docker Multi-stage Build
- Integration Tests
- Testcontainers
- SMS-провайдер
- Email-уведомления

---

# Используемые практики

- Clean Architecture (облегчённая)
- Repository Pattern
- Service Layer
- Dependency Injection
- JWT Authentication
- Transaction Management
- Idempotency
- Unit Testing
- Manual Testing
- Conventional Commits

---
