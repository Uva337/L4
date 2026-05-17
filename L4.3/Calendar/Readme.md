1. Запуск проекта
Для запуска сервера выполните команду в корне проекта:

go run cmd/app/main.go

По умолчанию сервер запускается на порту 8080.

Запуск тестов
Бизнес-логика покрыта unit-тестами. Для проверки выполните:

go test -v ./internal/service/...

1. Создание события (POST /create_event)
Запрос:

curl -X POST http://localhost:8080/create_event \
-H "Content-Type: application/json" \
-d '{
  "id": "evt_100",
  "user_id": 1,
  "title": "Собеседование на Golang разработчика",
  "date": "2026-05-20T14:00:00Z",
  "reminder_at": "2026-05-20T13:30:00Z"
}'
Успешный ответ: {"result": "событие создано"}

2. Получение событий за день (GET /events_for_day)
Запрос:


curl -X GET "http://localhost:8080/events_for_day?user_id=1"

Успешный ответ:
{
  "result": [
    {
      "id": "evt_100",
      "user_id": 1,
      "title": "Собеседование на Golang разработчика",
      "date": "2026-05-20T14:00:00Z",
      "reminder_at": "2026-05-20T13:30:00Z"
    }
  ]
}
3. Удаление события (POST /delete_event)
Запрос:

curl -X POST http://localhost:8080/delete_event \
-H "Content-Type: application/json" \
-d '{"id": "evt_100"}'

Успешный ответ: {"result": "событие удалено"}