package store

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrCodeAlreadyExists = errors.New("short code already exists")
	ErrCodeNotFound      = errors.New("short code not found")
)

type URLRecord struct {
	Code        string
	OriginalURL string
	CreatedAt   time.Time
	Visits      int
}

type URLStore interface {
	Save(ctx context.Context, record URLRecord) error
	GetByCode(ctx context.Context, code string) (URLRecord, error)
	IncrementVisitsAndGet(ctx context.Context, code string) (URLRecord, error)
	List(ctx context.Context) ([]URLRecord, error)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]URLRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]URLRecord),
	}
}

func (s *MemoryStore) Save(_ context.Context, record URLRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[record.Code]; exists {
		return ErrCodeAlreadyExists
	}

	s.data[record.Code] = record
	return nil
}

func (s *MemoryStore) GetByCode(_ context.Context, code string) (URLRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.data[code]
	if !exists {
		return URLRecord{}, ErrCodeNotFound
	}

	return record, nil
}

func (s *MemoryStore) IncrementVisitsAndGet(_ context.Context, code string) (URLRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.data[code]
	if !exists {
		return URLRecord{}, ErrCodeNotFound
	}

	record.Visits++
	s.data[code] = record
	return record, nil
}

func (s *MemoryStore) List(_ context.Context) ([]URLRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]URLRecord, 0, len(s.data))
	for _, record := range s.data {
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})

	return records, nil
}
