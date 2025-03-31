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

package service

import (
	"context"

	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
)

type UrlService interface {
	CreateShortUrl(
		ctx context.Context,
		userId string,
		req *valueobject.CreateUrlRequest,
	) (*valueobject.CreateUrlResponse, error)
	GetLongUrl(ctx context.Context, shortUrl string) (string, error)
	GetAnalytics(ctx context.Context, shortUrl string, userId string) (int, error)
	GetPaginatedUrls(
		ctx context.Context,
		userId string,
		limit int,
		offset int,
	) ([]valueobject.UrlResponse, error)
	UpdateUrl(ctx context.Context, urlId string, newUrl string) error
	DeleteUrl(ctx context.Context, urlId string, userId string) error
}
