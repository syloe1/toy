package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"url-shortener/internal/store"
)

var (
	ErrInvalidURL   = errors.New("invalid url")
	ErrCodeNotFound = errors.New("short code not found")
)

const (
	defaultCodeLength = 6
	maxCodeRetries    = 5
	codeAlphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type TimeNow func() time.Time

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type ShortURL struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
	Visits      int       `json:"visits"`
}

type URLService struct {
	store   store.URLStore
	baseURL string
	now     TimeNow
}

func NewURLService(store store.URLStore, baseURL string, now TimeNow) *URLService {
	if now == nil {
		now = time.Now
	}

	return &URLService{
		store:   store,
		baseURL: strings.TrimRight(baseURL, "/"),
		now:     now,
	}
}

func (s *URLService) Create(ctx context.Context, req CreateShortURLRequest) (ShortURL, error) {
	normalizedURL, err := normalizeURL(req.OriginalURL)
	if err != nil {
		return ShortURL{}, err
	}

	createdAt := s.now().UTC()
	for range maxCodeRetries {
		code, err := generateCode(defaultCodeLength)
		if err != nil {
			return ShortURL{}, fmt.Errorf("generate code: %w", err)
		}

		record := store.URLRecord{
			Code:        code,
			OriginalURL: normalizedURL,
			CreatedAt:   createdAt,
			Visits:      0,
		}

		if err := s.store.Save(ctx, record); err != nil {
			if errors.Is(err, store.ErrCodeAlreadyExists) {
				continue
			}
			return ShortURL{}, fmt.Errorf("save short url: %w", err)
		}

		return s.toShortURL(record), nil
	}

	return ShortURL{}, errors.New("could not generate unique short code")
}

func (s *URLService) Resolve(ctx context.Context, code string) (ShortURL, error) {
	record, err := s.store.IncrementVisitsAndGet(ctx, code)
	if err != nil {
		if errors.Is(err, store.ErrCodeNotFound) {
			return ShortURL{}, ErrCodeNotFound
		}
		return ShortURL{}, fmt.Errorf("resolve short code: %w", err)
	}

	return s.toShortURL(record), nil
}

func (s *URLService) GetByCode(ctx context.Context, code string) (ShortURL, error) {
	record, err := s.store.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, store.ErrCodeNotFound) {
			return ShortURL{}, ErrCodeNotFound
		}
		return ShortURL{}, fmt.Errorf("get short code: %w", err)
	}

	return s.toShortURL(record), nil
}

func (s *URLService) List(ctx context.Context) ([]ShortURL, error) {
	records, err := s.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list short urls: %w", err)
	}

	result := make([]ShortURL, 0, len(records))
	for _, record := range records {
		result = append(result, s.toShortURL(record))
	}

	return result, nil
}

func (s *URLService) toShortURL(record store.URLRecord) ShortURL {
	return ShortURL{
		Code:        record.Code,
		OriginalURL: record.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, record.Code),
		CreatedAt:   record.CreatedAt,
		Visits:      record.Visits,
	}
}

func normalizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrInvalidURL
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", ErrInvalidURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", ErrInvalidURL
	}

	if parsed.Host == "" {
		return "", ErrInvalidURL
	}

	return parsed.String(), nil
}

func generateCode(length int) (string, error) {
	var builder strings.Builder
	//提前分配好长度， 速度更快
	builder.Grow(length)
	//生成一个大整数
	limit := big.NewInt(int64(len(codeAlphabet)))
	for i := 0; i < length; i++ {
		//crypto/rand安全随机数生成器
		// 【0， limit-1]的数字
		//rand.Int生成一个随机大整数
		n, err := rand.Int(rand.Reader, limit)
		//生成随机数失败
		if err != nil {
			return "", err
		}
		//放入builder
		//n.Int64() 随机数
		builder.WriteByte(codeAlphabet[n.Int64()])
	}

	return builder.String(), nil
}
