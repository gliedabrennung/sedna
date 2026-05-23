# Sedna

Высоконагруженный бэкенд мессенджера реального времени.

## 🛠 Стек технологий
* **Язык**: Go 1.26
* **Web/API**: Hertz, WebSocket
* **Базы данных**: PostgreSQL (пользователи), ScyllaDB (история сообщений)
* **Кэш & Pub/Sub**: Redis
* **Инфраструктура**: Docker, Docker Compose

## 🚀 Запуск проекта

Убедитесь, что у вас установлены **Go 1.26+** и **Docker**.

1. Запустите инфраструктуру (БД и Redis):
   ```bash
   docker-compose up -d
   ```
2. Запустите сервер:
   ```bash
   go run cmd/server/main.go
   ```
3. Запуск тестов:
   ```bash
   go test ./..
   ```

## 📡 Примеры запросов

### 1. Авторизация (REST)
**Запрос:** `POST /api/v1/auth/login`
```json
{
  "username": "johndoe",
  "password": "secretpassword"
}
```
**Ответ:** `200 OK` (возвращает JWT токен)

### 2. Получение истории чата (REST)
**Запрос:** `GET /api/v1/messages/history?partner_id=2`
*(Не забудьте передать куку `token` или заголовок `Authorization: Bearer <token>`)*
**Ответ:** `200 OK`
```json
{
  "messages": [
    {
      "message_id": "uuid-v4",
      "content": "Привет!",
      "from_id": 1,
      "to_id": 2,
      "created_at": "2026-05-22T20:00:00Z"
    }
  ],
  "next_cursor": ""
}
```

### 3. Отправка сообщения (WebSocket)
**Подключение:** `ws://localhost:8080/ws`
**Payload (JSON):**
```json
{
  "to": 2,
  "message": "Привет! Как дела?"
}
```
