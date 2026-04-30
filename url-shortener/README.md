# Go URL Shortener

This project is a complete URL shortening web application written in Go.

## Architecture

The project follows the layered flow you requested:

`router -> handler -> service`

- `internal/router`: route registration and lightweight middleware
- `internal/handler`: HTTP request parsing and response shaping
- `internal/service`: business logic, validation, short-code generation
- `internal/store`: in-memory persistence abstraction

## Features

- Create short URLs from HTML form or JSON API
- Redirect short URLs to the original destination
- Query one short URL by code
- List all shortened URLs
- Track visit counts
- Basic health check endpoint
- Browser UI for quick manual testing

## Run

```bash
go run ./cmd/server
```

Open `http://localhost:8080`.

## API

### Health

```http
GET /api/health
```

### Create short URL

```http
POST /api/urls
Content-Type: application/json

{
  "original_url": "https://example.com/article/123"
}
```

### List short URLs

```http
GET /api/urls
```

### Get one short URL

```http
GET /api/urls/{code}
```

### Redirect

```http
GET /{code}
```
