package service

import (
	"context"
	"testing"
	"time"

	"url-shortener/internal/store"
)

func TestCreateAndResolve(t *testing.T) {
	repo := store.NewMemoryStore()
	now := func() time.Time {
		return time.Date(2026, 4, 30, 8, 0, 0, 0, time.UTC)
	}

	svc := NewURLService(repo, "http://localhost:8080", now)

	created, err := svc.Create(context.Background(), CreateShortURLRequest{
		OriginalURL: "https://example.com/docs",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.Code == "" {
		t.Fatal("Create() returned empty code")
	}

	if created.ShortURL == "" {
		t.Fatal("Create() returned empty short url")
	}

	resolved, err := svc.Resolve(context.Background(), created.Code)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if resolved.Visits != 1 {
		t.Fatalf("Resolve() visits = %d, want 1", resolved.Visits)
	}
}

func TestCreateRejectsInvalidURL(t *testing.T) {
	repo := store.NewMemoryStore()
	svc := NewURLService(repo, "http://localhost:8080", time.Now)

	_, err := svc.Create(context.Background(), CreateShortURLRequest{
		OriginalURL: "not-a-url",
	})
	if err == nil {
		t.Fatal("Create() expected error for invalid url")
	}

	if err != ErrInvalidURL {
		t.Fatalf("Create() error = %v, want %v", err, ErrInvalidURL)
	}
}
