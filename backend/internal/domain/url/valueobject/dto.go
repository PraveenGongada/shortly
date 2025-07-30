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

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CreateURLRequest represents URL creation request data
type CreateURLRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

// CreateURLResponse represents URL creation response data
type CreateURLResponse struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
}

// CreateShortURLResponse creates a CreateURLResponse from a URL entity
func CreateShortURLResponse(url *entity.URL) CreateURLResponse {
	return CreateURLResponse{
		ID:        url.ID(),
		ShortCode: url.ShortCode(),
	}
}

// URLResponse represents URL data in responses
type URLResponse struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
	Redirects int    `json:"redirects"`
}

// CreateGetURLsResponse creates a slice of URLResponse from URL entities
func CreateGetURLsResponse(urls []*entity.URL) []URLResponse {
	urlResponse := make([]URLResponse, len(urls))

	for i, url := range urls {
		urlResponse[i] = URLResponse{
			ID:        url.ID(),
			ShortCode: url.ShortCode(),
			LongURL:   url.LongURL(),
			Redirects: url.Redirects(),
		}
	}

	return urlResponse
}

// URLUpdateRequest represents URL update request data
type URLUpdateRequest struct {
	ID     string `json:"id" validate:"required"`
	NewURL string `json:"new_url" validate:"required,url"`
}

// DeleteURLRequest represents URL deletion request data
type DeleteURLRequest struct {
	ID string `json:"id" validate:"required"`
}
