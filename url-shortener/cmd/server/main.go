package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"url-shortener/internal/handler"
	"url-shortener/internal/router"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
)

func main() {
	addr := getEnv("APP_ADDR", ":8080")
	baseURL := getEnv("APP_BASE_URL", "http://localhost:8080")

	repo := store.NewMemoryStore()
	urlService := service.NewURLService(repo, baseURL, time.Now)
	urlHandler := handler.NewURLHandler(urlService)
	appRouter := router.New(urlHandler)

	log.Printf("url-shortener is listening on %s", addr)
	if err := http.ListenAndServe(addr, appRouter); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
