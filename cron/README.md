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
