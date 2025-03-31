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

package valueobject

import "github.com/PraveenGongada/shortly/internal/domain/url/entity"

type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type CreateUrlRequest struct {
	LongUrl string `json:"long_url"`
}

type CreateUrlResponse struct {
	Id       string `json:"id"`
	ShortUrl string `json:"short_url"`
}

func CreateShortUrlResponse(url *entity.Url) CreateUrlResponse {
	return CreateUrlResponse{
		Id:       url.Id,
		ShortUrl: url.ShortUrl,
	}
}

type UrlResponse struct {
	Id        string `json:"id"`
	ShortUrl  string `json:"shortUrl"`
	LongUrl   string `json:"longUrl"`
	Redirects int    `json:"redirects"`
}

func CreateGetUrlsResponse(urls []entity.Url) []UrlResponse {
	urlResponse := make([]UrlResponse, len(urls))

	for i, url := range urls {
		urlResponse[i] = UrlResponse{
			Id:        url.Id,
			ShortUrl:  url.ShortUrl,
			LongUrl:   url.LongUrl,
			Redirects: url.Redirects,
		}
	}

	return urlResponse
}

type UrlUpdateRequest struct {
	Id  string `json:"id"`
	Url string `json:"new_url"`
}

type DeleteUrlRequest struct {
	Id string `json:"id"`
}
