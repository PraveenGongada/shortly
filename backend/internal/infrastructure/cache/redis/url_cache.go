package redis

import (
	"context"
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/url/cache"
)

type urlCache struct {
	client Client
}

func NewURLCache(client Client) cache.URLCache {
	return &urlCache{
		client: client,
	}
}

func (uc *urlCache) SetShortURL(
	ctx context.Context,
	shortCode, originalURL string,
	ttl time.Duration,
) error {
	return nil
}

func (uc *urlCache) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return "", nil
}

func (uc *urlCache) InvalidateShortURL(ctx context.Context, shortCode string) error {
	return nil
}
