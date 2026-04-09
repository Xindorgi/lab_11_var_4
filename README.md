# Лабораторная работа 11 — Вариант 4 — Ведешкин Андрей Георгиевич — группа 221131

---

## Общая структура проекта

```text
./
├── compose.yaml               # Оркестрация контейнеров
├── go-app/
│   ├── main.go                # HTTP-сервер на Go (порт 8080)
│   ├── main_test.go           # Юнит-тесты Go
│   ├── go.mod
│   └── Dockerfile             # Многоэтапная сборка, финальный образ scratch
└── python-app/
    ├── app.py                 # HTTP-сервер на FastAPI (порт 8000)
    ├── test_app.py            # Юнит-тесты Python
    ├── requirements.txt
    └── Dockerfile             # Образ python:3.12-alpine
```

---

## Реализованные задания

- **M4.** Собрать образы и сравнить их размеры
- **M6.** Настроить сеть между контейнерами
- **M8.** Добавить healthcheck для каждого сервиса
- **H4.** Использовать docker buildx для кросс-платформенной сборки (arm64/amd64)
- **H6.** Реализовать graceful shutdown и проброс сигналов в контейнеры

---

## M4 — Сравнение размеров образов

Для Go используется многоэтапная сборка: на первом этапе компилируется статически слинкованный бинарник (`CGO_ENABLED=0`), на втором — он копируется в пустой образ `FROM scratch`. Результат — образ содержит только исполняемый файл без ОС, рантайма и пакетного менеджера.

Для Python Dockerfile поддерживает переопределение базового образа через `ARG BASE`, что позволяет сравнить alpine и slim варианты одной командой.

**Сборка и сравнение:**

```bash
# Go (scratch)
docker build -t go-app:scratch ./go-app

# Python alpine (по умолчанию)
docker build -t python-app:alpine ./python-app

# Python slim
docker build -t python-app:slim --build-arg BASE=python:3.12-slim ./python-app

# Сравнение размеров
docker images go-app python-app
```

**Результаты:**

| Образ | Базовый образ | Размер |
|-------|--------------|--------|
| go-app:scratch | scratch | ~7.82 MB |
| python-app:alpine | python:3.12-alpine | ~184 MB |
| python-app:slim | python:3.12-slim | ~300 MB |

Go-образ в ~23 раза меньше Python alpine, так как содержит только один статический бинарник. Python несёт в себе интерпретатор, стандартную библиотеку и pip-зависимости.

---

## M6 — Сеть между контейнерами

Оба сервиса объединены в пользовательскую bridge-сеть `app-net` в `compose.yaml`. Docker автоматически обеспечивает DNS-резолвинг по имени сервиса — python-app обращается к go-app как `http://go-app:8080`.

Переменная окружения `GO_APP_URL=http://go-app:8080/time` передаётся в python-app. Эндпоинт `/check-go` в Python использует её для запроса к Go-сервису и возвращает ответ — так проверяется, что сеть работает.

**Проверка:**

```bash
docker compose up -d

# Убедиться, что python-app достучался до go-app
curl http://localhost:8000/check-go
```

Ожидаемый ответ содержит текущее время от Go-сервиса: `{"time": "2026-04-09T16:26:15Z"}`.

---

## M8 — Healthcheck

Каждый сервис имеет healthcheck как в Dockerfile, так и в compose.yaml.

**Go** использует встроенный флаг `-health`: приложение делает HTTP-запрос к собственному `/health` и завершается с кодом `0` (успех) или `1` (ошибка). Это позволяет работать из образа `scratch`, где нет `curl` или `wget`.

```dockerfile
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "-health"]
```

**Python** использует `curl`:

```dockerfile
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD curl --fail http://localhost:8000/health || exit 1
```

В `compose.yaml` прописан `condition: service_healthy` — python-app запускается только после того, как go-app перешёл в статус `healthy`.

**Проверка статусов:**

```bash
docker compose up -d
docker compose ps
# Оба контейнера должны показывать статус healthy
```

---

## H4 — Кросс-платформенная сборка (buildx)

`go-app/Dockerfile` принимает `ARG TARGETOS` и `ARG TARGETARCH`, которые Docker buildx подставляет автоматически при сборке под целевую платформу. Go компилирует бинарник под нужную ОС и архитектуру через `GOOS`/`GOARCH`.

**Настройка buildx:**

```bash
docker buildx create --name multiarch --use
docker buildx inspect --bootstrap
```

**Сборка multi-arch образа с пушем в registry:**

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t your-registry/go-app:latest \
  --push \
  ./go-app
```

**Локальная сборка под конкретную архитектуру (без registry):**

```bash
# amd64
docker buildx build --platform linux/amd64 -t go-app:amd64 --load ./go-app

# arm64
docker buildx build --platform linux/arm64 -t go-app:arm64 --load ./go-app
```

**Проверка архитектуры образа:**

```bash
docker inspect go-app:amd64 | grep Architecture
docker inspect go-app:arm64 | grep Architecture
```

---

## H6 — Graceful Shutdown

Оба приложения корректно обрабатывают сигнал `SIGTERM`, который Docker отправляет при `docker compose stop` или `docker stop`.

**Go** — функция `runServer` подписывается на `SIGTERM` и `SIGINT` через `signal.Notify`. При получении сигнала вызывается `srv.Shutdown` с таймаутом 5 секунд: сервер перестаёт принимать новые соединения и ждёт завершения активных запросов.

```go
case sig := <-quit:
    log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
        return err
    }
    log.Println("Server gracefully stopped")
```

**Python** — класс `Server` регистрирует обработчик `SIGTERM`, который устанавливает флаг `server.should_exit = True`. Uvicorn при следующей итерации event loop корректно завершает все соединения.

В `compose.yaml` добавлен `stop_grace_period: 10s` — Docker будет ждать 10 секунд до принудительного `SIGKILL`.

Важно: в Dockerfile используется exec-форма `ENTRYPOINT ["./app"]`, а не shell-форма. В shell-форме PID 1 — это `/bin/sh`, и сигнал до приложения не доходит.

**Демонстрация:**

```bash
docker compose up -d
docker compose stop
docker compose logs
```

Ожидаемые логи:

```
python-app  | INFO:     Shutting down
python-app  | INFO:     Waiting for application shutdown.
python-app  | INFO:     Application shutdown complete.
go-app      | 2026/04/09 16:26:49 Received signal: terminated, initiating graceful shutdown...
go-app      | 2026/04/09 16:26:49 Server gracefully stopped
python-app  | INFO:     Finished server process [1]
```

---

## Запуск

### Запуск через Docker Compose

```bash
docker compose up --build -d
```

### Проверка эндпоинтов

```bash
# Go
curl http://localhost:8080/time
curl http://localhost:8080/health

# Python
curl http://localhost:8000/ping
curl http://localhost:8000/health
curl http://localhost:8000/check-go
```

### Запуск юнит-тестов

```bash
# Go
cd go-app
go test ./... -v

# Python
cd python-app
pip install -r requirements.txt
pytest test_app.py -v
```

### Остановка

```bash
docker compose stop
```

---

## ✅ Статус тестов

- ✅ Юнит-тесты Go — 5/5 passed
- ✅ Юнит-тесты Python — 6/6 passed
- ✅ Healthcheck Go (scratch, `-health` флаг)
- ✅ Healthcheck Python (curl)
- ✅ Сеть между контейнерами (app-net, DNS по имени сервиса)
- ✅ Кросс-платформенная сборка (amd64/arm64)
- ✅ Graceful shutdown с явным логированием сигналов

**Все тесты проходят успешно!**