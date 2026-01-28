# Chats API Service
REST API сервис для чата, построенный на Go с использованием PostgreSQL.

### Технологии:
- Backend: Go 1.25 (net/http)
- Database: PostgreSQL
- ORM: GORM
- Router: Chi
- Containerization: Docker & Docker Compose
- Migrations: Goose
- Testing: testify

### Архитектура:
- domain - сущности/модели
- repositories - работа с базой данных
- services - бизнес-логика
- handlers - HTTP обработчики (тесты)
- config - конфигурация
- database - подключение к БД
- route - маршруты
- helpers - вспомогательные функции
- migrations - миграции

### Запуск сервиса

1. Клонируйте репозиторий

`git clone <repository-url>`

2. Перейдите в папку с проектом

`cd api_service_chat`

3. Создайте и запустите контейнер

`docker-compose build`
`docker-compose up -d`

Сервис будет доступен по адресу:

`http://localhost:8080`

4. Проверьте работу:

`curl http://localhost:8080/health`

5. Чтобы посмотреть логи контейнера, выполните команду:

`docker-compose logs -f app`

6. Чтобы остановить и удалить контейнер, выполните команду:

`docker-compose down`

7. Запуск тестов:

`go test ./internal/handlers -v`

## API Endpoints

### Chats:

- GET `/api/chats/{id}` — получить чат и последние N сообщений
- POST `/api/chats` — создать новый чат
- DELETE `/api/chats/{id}` — удалить чат вместе со всеми сообщениями

### Messages:

- POST `/api/chats/{id}/messages` — отправить сообщение в чат

### Health Check:

- GET `/health` - проверка статуса API
