# Calculator

## Архитектура
```mermaid
graph TD
  Orchestrator -->|HTTP| Client
  Agent2 -->|gRPC| Orchestrator
  Agent1 -->|gRPC| Orchestrator
  Orchestrator -->|gRPC| Agent2[Агент 2]
  Orchestrator -->|gRPC| Agent1[Агент 1]
  Client[Клиент] -->|HTTP| Orchestrator[Оркестратор]
```
## Установка и запуск
1. Клонирование репозитория

```bash
git clone https://github.com/Dided08/Calculator.git
cd Calculator
```
   

3. Требования

    Go 1.20+

    Docker и Docker Compose (опционально)

4. Запуск оркестратора
   
Linux/macOS
```bash
export TIME_ADDITION_MS=200
export TIME_SUBTRACTION_MS=200
export TIME_MULTIPLICATIONS_MS=300
export TIME_DIVISIONS_MS=400
   
go run ./cmd/orchestrator/main.go
```

Windows PowerShell
```bash
$env:TIME_ADDITION_MS=200
$env:TIME_SUBTRACTION_MS=200
$env:TIME_MULTIPLICATIONS_MS=300
$env:TIME_DIVISIONS_MS=400
go run .\cmd\orchestrator\main.go
```

4. Запуск агента
   
Linux/macOS
```bash
export COMPUTING_POWER=4
export ORCHESTRATOR_URL="localhost:8080"
go run ./cmd/agent/main.go
```

Windows PowerShell
```bash
$env:COMPUTING_POWER=4
$env:ORCHESTRATOR_URL="localhost:8080"
go run .\cmd\agent\main.go
```
5. Запуск фронтенда

В браузере перейдите по ссылке http://localhost:8081/

6. Docker Compose (опционально)
```bash
docker-compose up --build
```

Если выводит project name must not be empty, то пропишите:
```bash
docker build -f Dockerfile.orchestrator -t orchestrator .
docker-compose up --build
```
   

## API (REST)

### POST /api/v1/calculate

Запускает вычисление выражения.

Запрос
```http
POST /api/v1/calculate HTTP/1.1
Content-Type: application/json
Authorization: Bearer <token>

{"expression":"(2+3)*4-10/2"}
```
    
Ответ (201 Created)
```json
{"id": 1}
```

### GET /api/v1/expressions

Возвращает все выражения пользователя.

Запрос
```http
GET /api/v1/expressions HTTP/1.1
Authorization: Bearer <token>
```

Ответ (200 OK)
   ```json
{
  "expressions": [
     {"id":1, "expression":"(2+3)*4-10/2", "status":"done", "result":15}
  ]
}
   ```

### GET /api/v1/expressions/{id}

Возвращает статус и результат выражения по его ID.

 Запрос
```http
GET /api/v1/expressions/1 HTTP/1.1
Authorization: Bearer <token>
```

Ответ (200 OK)
```json
{"expression": {"id":1, "status":"done", "result":15}}
```

## Примеры использования

### Простое выражение
```bash
    curl -X POST http://localhost:8080/api/v1/calculate \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"expression":"((3+5)*2-8)/4"}'
# {"id":1}
    
    curl http://localhost:8080/api/v1/expressions/1 \
      -H "Authorization: Bearer $TOKEN"
# {"expression":{"id":1,"status":"done","result":2}}
```

### Ошибка деления на ноль
```bash
    curl -X POST http://localhost:8080/api/v1/calculate \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{"expression":"10/(5-5)"}'
# HTTP 422: invalid expression or result out of range
```

### Тестирование
```bash
go test -v ./cmd/agent
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
| --- | --- | --- |
| TIME_ADDITION_MS | Задержка для операции + (в миллисекундах) | 100 |
| TIME_SUBTRACTION_MS | Задержка для операции - | 100 |
| TIME_MULTIPLICATIONS_MS | Задержка для операции * | 100 |
| TIME_DIVISIONS_MS | Задержка для операции / | 100 |
| COMPUTING_POWER | Количество потоков обработки у агента | 100 |
| ORCHESTRATOR_URL	 | Адрес gRPC-оркестратора (например, host:port) | 100 |
