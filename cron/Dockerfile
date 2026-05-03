FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget
RUN addgroup -S app && adduser -S app -G app

ENV GIN_MODE=release
ENV CONFIG_PATH=/app/config.yaml

COPY --from=builder /out/server /app/server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/healthz >/dev/null || exit 1

USER app

CMD ["/app/server"]
