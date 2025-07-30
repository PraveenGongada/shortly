package cache

import (
	"context"
	"time"
)

type URLCache interface {
	SetShortURL(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	InvalidateShortURL(ctx context.Context, shortCode string) error
}
