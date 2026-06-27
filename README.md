# Rocket Factory

`Rocket Factory` — учебный микросервисный проект на Go, в котором заказ собирается из нескольких сервисов и протоколов.

Проект показывает типичную схему backend-монорепозитория:

- внешний REST API для клиентского сервиса;
- внутреннее gRPC-взаимодействие между микросервисами;
- отдельные хранилища под разные bounded context;
- общие платформенные пакеты;
- генерация кода из OpenAPI и Protobuf;
- контейнерный запуск и CI/CD через GitHub Actions.

## Что делает проект

В репозитории реализованы три бизнес-сервиса:

- `order` — сервис заказов. Принимает HTTP-запросы, создает заказ, читает заказ, отменяет его и инициирует оплату.
- `inventory` — сервис склада/каталога деталей. Отдает детали по gRPC, умеет искать детали по фильтру.
- `payment` — сервис оплаты. Сейчас это легковесный payment stub: он валидирует способ оплаты и возвращает `transaction_uuid`.

Сценарий работы выглядит так:

1. Клиент вызывает `order` по HTTP.
2. `order` по gRPC обращается в `inventory`, чтобы проверить и получить детали.
3. `order` сохраняет заказ в PostgreSQL.
4. При оплате `order` по gRPC вызывает `payment`.
5. `payment` возвращает `transaction_uuid`, после чего `order` обновляет статус заказа.

## Архитектура

```text
Client
  |
  | HTTP / OpenAPI
  v
order-service
  | \
  |  \ gRPC
  |   \
  v    v
inventory-service   payment-service
      |                  |
      v                  v
   MongoDB            stateless logic

order-service -> PostgreSQL
```

Ключевая идея проекта:

- наружу торчит только `order` как REST API;
- внутренние интеграции между сервисами построены на gRPC;
- контракты лежат в `shared`, а общая инфраструктурная логика — в `platform`.

## Модули монорепозитория

Проект организован как Go workspace через [`go.work`](./go.work).

В workspace подключены отдельные модули:

- `order`
- `inventory`
- `payment`
- `platform`
- `shared`

Это значит, что каждый сервис имеет свой `go.mod`, но локально они развиваются как единая система.

## Структура репозитория

```text
.
├── order/              # REST API заказов + PostgreSQL + gRPC adapters
├── inventory/          # gRPC API деталей + MongoDB
├── payment/            # gRPC API оплаты
├── platform/           # общие инфраструктурные пакеты
├── shared/             # OpenAPI / Proto контракты и сгенерированный код
├── deploy/
│   ├── compose/        # docker compose по сервисам
│   ├── docker/         # Dockerfile'ы
│   └── env/            # env-файлы для compose
├── Taskfile.yml        # основные команды проекта
└── .github/workflows/  # CI
```

## Сервисы

### Order Service

`order` — точка входа в систему.

Что делает:

- принимает HTTP-запросы на создание и управление заказами;
- обращается в `inventory` по gRPC для получения деталей;
- обращается в `payment` по gRPC для оплаты;
- хранит заказы в PostgreSQL;
- использует миграции через `goose`;
- генерирует REST-слой из OpenAPI через `ogen`.

Основные каталоги:

- [`order/cmd/order-service`](./order/cmd/order-service)
- [`order/internal/app`](./order/internal/app)
- [`order/internal/api`](./order/internal/api)
- [`order/internal/service`](./order/internal/service)
- [`order/internal/repository/postgres`](./order/internal/repository/postgres)
- [`order/internal/integration/grpc`](./order/internal/integration/grpc)
- [`order/migrations`](./order/migrations)

### Inventory Service

`inventory` — внутренний gRPC-сервис деталей.

Что делает:

- отдает детали по UUID;
- отдает список деталей по фильтрам;
- хранит данные в MongoDB;
- может инициализировать данные фикстурами;
- использует protobuf/gRPC контракты из `shared`.

Основные каталоги:

- [`inventory/cmd/inventory-service`](./inventory/cmd/inventory-service)
- [`inventory/internal/app`](./inventory/internal/app)
- [`inventory/internal/api`](./inventory/internal/api)
- [`inventory/internal/service`](./inventory/internal/service)
- [`inventory/internal/repository/mongo`](./inventory/internal/repository/mongo)

### Payment Service

`payment` — внутренний gRPC-сервис оплаты.

Что делает:

- принимает запрос на оплату заказа;
- проверяет валидность `payment_method`;
- возвращает `transaction_uuid`;
- не использует внешнюю платежную систему и не хранит состояние в БД.

Сервис сейчас специально простой — это хороший stub для отладки оркестрации заказов.

Основные каталоги:

- [`payment/cmd/payment-service`](./payment/cmd/payment-service)
- [`payment/internal/app`](./payment/internal/app)
- [`payment/internal/api`](./payment/internal/api)
- [`payment/internal/service`](./payment/internal/service)

## Общие модули

### platform

[`platform`](./platform) — внутренний технический слой, который переиспользуют сервисы.

Содержит:

- структурированный логгер на `slog`;
- `closer` для контролируемого shutdown;
- gRPC interceptors;
- HTTP middleware;
- общие ошибки;
- helper для UUID parsing;
- gRPC health registration.

Ключевые пакеты:

- [`platform/pkg/logger`](./platform/pkg/logger)
- [`platform/pkg/closer`](./platform/pkg/closer)
- [`platform/pkg/interceptor`](./platform/pkg/interceptor)
- [`platform/pkg/middleware`](./platform/pkg/middleware)
- [`platform/pkg/error`](./platform/pkg/error)
- [`platform/pkg/uuidx`](./platform/pkg/uuidx)

### shared

[`shared`](./shared) — место для контрактов и сгенерированного кода.

Содержит:

- OpenAPI спецификацию `order` API;
- Protobuf схемы для `inventory` и `payment`;
- сгенерированный Go-код для REST и gRPC.

Ключевые каталоги:

- [`shared/api/order/v1`](./shared/api/order/v1)
- [`shared/proto`](./shared/proto)
- [`shared/pkg/openapi/order/v1`](./shared/pkg/openapi/order/v1)
- [`shared/pkg/proto`](./shared/pkg/proto)

## Контракты API

### REST: order-service

OpenAPI-описание находится в [`shared/api/order/v1/order.openapi.yaml`](./shared/api/order/v1/order.openapi.yaml).

Основные маршруты:

- `POST /api/v1/orders`
- `GET /api/v1/orders/{order_uuid}`
- `POST /api/v1/orders/{order_uuid}/pay`
- `POST /api/v1/orders/{order_uuid}/cancel`

### gRPC: inventory-service

Proto-файл: [`shared/proto/inventory/v1/inventory.proto`](./shared/proto/inventory/v1/inventory.proto)

Методы:

- `GetPart`
- `ListParts`

### gRPC: payment-service

Proto-файл: [`shared/proto/payment/v1/payment.proto`](./shared/proto/payment/v1/payment.proto)

Методы:

- `PayOrder`

## Технологии и инструменты

### Язык и платформа

- Go `1.25.6`
- Go Workspace (`go.work`)
- стандартный логгер `log/slog`

### Транспорт и контракты

- HTTP
- REST
- OpenAPI 3.0.3
- gRPC
- Protocol Buffers
- `ogen` для генерации REST-кода
- `buf` + `protoc-gen-go` + `protoc-gen-go-grpc` для генерации protobuf/gRPC-кода

### Web и серверная часть

- `chi` для HTTP router/middleware
- `google.golang.org/grpc`

### Хранилища

- PostgreSQL для `order`
- MongoDB для `inventory`
- `payment` сейчас без отдельной БД

### Работа с БД

- `pgx/v5` и `scany`
- `squirrel` для SQL builder
- `goose` для миграций PostgreSQL
- официальный MongoDB Go driver

### Тесты и качество кода

- `testify`
- `mockery`
- `golangci-lint`
- `gofumpt`
- `gci`
- `gosec`

### DevOps и инфраструктура

- Docker
- Docker Compose
- GitHub Actions
- `go-task` / `Taskfile`
- `grpcurl` для smoke/API тестов

## Конфигурация

Примеры окружения лежат в:

- [`deploy/env/order/.env`](./deploy/env/order/.env)
- [`deploy/env/inventory/.env`](./deploy/env/inventory/.env)
- [`deploy/env/payment/.env`](./deploy/env/payment/.env)

### Order

Основные переменные:

- `DB_USER`
- `DB_PASSWORD`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `ORDER_SERVICE_HOST`
- `ORDER_SERVICE_PORT`
- `INVENTORY_ADAPTER_HOST`
- `INVENTORY_ADAPTER_PORT`
- `PAYMENT_ADAPTER_HOST`
- `PAYMENT_ADAPTER_PORT`
- `LOGGER_LEVEL`
- `LOGGER_AS_JSON`

### Inventory

Основные переменные:

- `DB_USER`
- `DB_PASSWORD`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_AUTH`
- `DB_IS_INIT`
- `SERVICE_HOST`
- `SERVICE_PORT`
- `LOGGER_LEVEL`
- `LOGGER_AS_JSON`

### Payment

Основные переменные:

- `SERVICE_HOST`
- `SERVICE_PORT`
- `LOGGER_LEVEL`
- `LOGGER_AS_JSON`

## Запуск проекта

### Что нужно установить

- Go `1.25.6`
- Docker + Docker Compose
- `task` CLI

Установка `task`:

Ubuntu:

```bash
sudo snap install task --classic
```

macOS:

```bash
brew install go-task
```

### Быстрый старт через Docker

Поднять всю систему:

```bash
task containers-up
```

Остановить всю систему:

```bash
task containers-down
```

Очистить volume'ы БД:

```bash
task clear-volumes
```

### Поднять только один сервис

```bash
task core-up
task inventory-service-up
task payment-service-up
task order-service-up
```

Или по отдельности остановить:

```bash
task inventory-service-down
task payment-service-down
task order-service-down
task core-down
```

### Локальный запуск без Docker

Технически сервисы можно запускать напрямую:

```bash
go run ./inventory/cmd/inventory-service
go run ./payment/cmd/payment-service
go run ./order/cmd/order-service
```

Но для этого нужно заранее подготовить:

- env-переменные;
- PostgreSQL для `order`;
- MongoDB для `inventory`;
- адреса gRPC-адаптеров для `order`.

## Полезные команды

### Форматирование

```bash
task format
```

### Линтинг

```bash
task lint
```

### Тесты

```bash
task test
task test-coverage
task test-coverage-html
```

### Генерация кода

OpenAPI:

```bash
task gen-ogen
```

Protobuf / gRPC:

```bash
task lint-proto
task gen-proto
```

Mocks:

```bash
task gen-mockery
```

Обновление зависимостей:

```bash
task update-deps
```

### Smoke/API тест проекта

```bash
task test-api
```

Эта задача:

- проверяет gRPC методы `inventory`;
- создает заказ через REST;
- оплачивает заказ;
- проверяет смену статуса;
- создает и отменяет второй заказ.

## Базы данных и хранение состояния

### PostgreSQL в order-service

Миграции лежат в [`order/migrations`](./order/migrations).

На данный момент там описаны:

- enum `payment_method`;
- enum `order_status`;
- таблица `orders`.

### MongoDB в inventory-service

Mongo используется как хранилище деталей.

Особенности:

- создается индекс для коллекции деталей;
- сервис может инициализировать набор тестовых деталей;
- фикстуры генерируются внутри проекта.

### Payment-service

`payment` пока не хранит платежи в отдельной БД и работает как stub.

Это упрощает разработку потока:

- создать заказ;
- вызвать оплату;
- получить `transaction_uuid`;
- обновить статус заказа.

## Внутренний паттерн сервисов

Во всех сервисах проект постепенно выровнен под единый bootstrap-паттерн:

- `cmd/<service>/main.go` — тонкая точка входа;
- `internal/app` — инициализация зависимостей и lifecycle сервиса;
- `internal/api` — транспортный слой;
- `internal/service` — бизнес-логика;
- `internal/repository` — работа с данными;
- `internal/config` — env-конфигурация.

Это делает проект более предсказуемым и упрощает сопровождение.

## CI/CD

GitHub Actions лежат в [`/.github/workflows`](./.github/workflows).

Текущие workflow:

- `ci.yaml` — общий entrypoint для CI;
- `lint.yaml` — линтинг всех модулей;
- `test.yaml` — запуск тестов всех модулей.

CI использует:

- `actions/setup-go`
- `go-task/setup-task`
- команды из `Taskfile.yml`

## Качество кода и стандарты

В проекте используются:

- строгий `golangci-lint` конфиг из [`.golangci.yml`](./.golangci.yml);
- запрет на `fmt.Print*` для логирования;
- запрет устаревшего `io/ioutil`;
- проверка безопасности через `gosec`;
- автоформатирование `gofumpt` + `gci`;
- генерация моков через [`.mockery.yaml`](./.mockery.yaml).

## Генерация контрактов

### OpenAPI / ogen

REST API `order` описан в OpenAPI, а затем генерируется в Go через `ogen`.

Конфиги:

- [`shared/api/order/v1/order.openapi.yaml`](./shared/api/order/v1/order.openapi.yaml)
- [`ogen.yaml`](./ogen.yaml)

### Protobuf / buf

gRPC контракты описаны в `shared/proto`.

Конфиги:

- [`shared/proto/buf.yaml`](./shared/proto/buf.yaml)
- [`shared/proto/buf.gen.yaml`](./shared/proto/buf.gen.yaml)

## Docker и деплой

Compose-файлы разбиты по сервисам:

- [`deploy/compose/core/compose.yaml`](./deploy/compose/core/compose.yaml)
- [`deploy/compose/order/compose.yaml`](./deploy/compose/order/compose.yaml)
- [`deploy/compose/inventory/compose.yaml`](./deploy/compose/inventory/compose.yaml)
- [`deploy/compose/payment/compose.yaml`](./deploy/compose/payment/compose.yaml)

Dockerfile'ы:

- [`deploy/docker/order/Dockerfile`](./deploy/docker/order/Dockerfile)
- [`deploy/docker/order/migrate.Dockerfile`](./deploy/docker/order/migrate.Dockerfile)
- [`deploy/docker/inventory/Dockerfile`](./deploy/docker/inventory/Dockerfile)
- [`deploy/docker/payment/Dockerfile`](./deploy/docker/payment/Dockerfile)

В контейнерах используется multi-stage build и запускаются только собранные бинарники.

## На что стоит обратить внимание

- `order` — orchestration service, именно он связывает остальные микросервисы.
- `inventory` и `payment` не торчат наружу REST API и работают как внутренние gRPC сервисы.
- `shared` и `platform` — важная часть архитектуры, потому что они выносят контракты и общую инфраструктуру из бизнес-сервисов.
- проект хорошо подходит как база для дальнейшего роста: auth, tracing, message broker, idempotency, saga/outbox, observability.

## Текущее состояние проекта

Сейчас репозиторий уже покрывает хороший базовый сценарий микросервисной системы:

- HTTP gateway/service facade для заказов;
- внутренние gRPC интеграции;
- разные базы под разные сервисы;
- генерация контрактов;
- миграции;
- контейнеризация;
- автоматизация через Taskfile;
- CI на GitHub Actions.

При этом проект остается достаточно компактным, чтобы его можно было развивать как учебный, pet-проект или основу для собеседовательного портфолио.

## Примечание

В `Taskfile.yml` smoke-тест `test-api` обращается к `http://127.0.0.1:8080`.

Проверь, чтобы внешний порт `order-service` в compose/env соответствовал этому адресу, если будешь запускать тест без дополнительных правок переменных окружения.
