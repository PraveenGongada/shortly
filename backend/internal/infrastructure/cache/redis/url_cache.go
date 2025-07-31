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
