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

package repository

import (
	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
)

type UrlRepository interface {
	CreateShortUrl(
		id string,
		userId string,
		shortUrl string,
		req *valueobject.CreateUrlRequest,
	) (*entity.Url, error)
	GetLongUrl(shortUrl string) (string, error)
	GetAnalytics(shortUrl string, userId string) (int, error)
	GetPaginatedUrls(userId string, limit int, offset int) ([]entity.Url, error)
	UpdateUrl(urlId string, newUrl string) error
	DeleteUrl(urlId string, userId string) error
}
