# Go TodoList

A single-user TodoList backend service built with Go, Gin, GORM, MySQL, and cron.

## Features

- Create, update, delete, and list tasks
- Batch delete tasks by ID list
- Update task status between `pending` and `done`
- Send daily reminder emails for unfinished tasks at 8:00 AM
- Clean `router -> handler -> service -> dao` layered architecture

## Project Structure

```text
cmd/server
internal/config
internal/dao
internal/handler
internal/mailer
internal/model
internal/router
internal/scheduler
internal/service
pkg/response
```

## Quick Start

1. Edit `config.yaml` with your MySQL DSN and SMTP settings.
2. Create the MySQL database, for example `todolist`.
3. Install dependencies:

```bash
go mod tidy
```

4. Run the server:

```bash
go run ./cmd/server
```

## Docker Deploy

Build and start the service with Docker Compose:

```bash
docker compose up -d --build
```

The compose file starts both the Go service and MySQL. Before deploying, replace the placeholder SMTP settings in `docker-compose.yml`.

If you already have a managed MySQL instance on your cloud server or cloud provider, keep only the `app` service and point `DATABASE_DSN` to that external database.

The application also supports environment variable overrides, which is convenient for container deployment:

- `SERVER_PORT`
- `DATABASE_DSN`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_FROM`
- `REMINDER_RECIPIENTS` (comma-separated)
- `REMINDER_CRON_SPEC`
- `REMINDER_TIMEZONE`

## API

### Health Check

- `GET /healthz`

### Task APIs

- `POST /api/v1/tasks`
- `GET /api/v1/tasks?status=all|pending|done`
- `PUT /api/v1/tasks/:id`
- `PATCH /api/v1/tasks/:id/status`
- `DELETE /api/v1/tasks/:id`
- `DELETE /api/v1/tasks`

### Request Samples

Create task:

```json
{
  "title": "Buy groceries",
  "content": "Milk, eggs, bread"
}
```

Update task status:

```json
{
  "status": "done"
}
```

Batch delete tasks:

```json
{
  "ids": [1, 2, 3]
}
```

## Notes

- New tasks default to `pending`.
- Reminder emails are sent only when unfinished tasks exist.
- SMTP send failures are logged and do not stop the service.
