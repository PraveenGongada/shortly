/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/url/cache"
)

type urlCache struct {
	client Client
	logger logger.Logger
}

func NewURLCache(client Client, logger logger.Logger) cache.URLCache {
	return &urlCache{
		client: client,
		logger: logger,
	}
}

func (uc *urlCache) SetShortURL(
	ctx context.Context,
	shortCode, originalURL string,
	ttl time.Duration,
) error {
	err := uc.client.Client().Set(ctx, shortCode, originalURL, ttl).Err()
	if err != nil {
		uc.logger.Error(ctx, "Error setting shortURL in cache",
			logger.String("shortCode", shortCode),
			logger.String("operation", "SetShortURL"),
			logger.Error(err),
		)
		return err
	}

	return nil
}

func (uc *urlCache) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := uc.client.Client().Get(ctx, shortCode).Result()
	if err != nil {
		uc.logger.Error(ctx, "Error getting shortURL in cache",
			logger.String("shortCode", shortCode),
			logger.String("operation", "GetOriginalURL"),
			logger.Error(err),
		)
		return "", err
	}

	return url, nil
}

func (uc *urlCache) InvalidateShortURL(ctx context.Context, shortCode string) error {
	err := uc.client.Client().Del(ctx, shortCode).Err()
	if err != redis.Nil && err != nil {
		uc.logger.Error(ctx, "Error invalidating shortURL in cache",
			logger.String("shortCode", shortCode),
			logger.String("operation", "InvalidateShortURL"),
			logger.Error(err),
		)
		return err
	}

	return nil
}
