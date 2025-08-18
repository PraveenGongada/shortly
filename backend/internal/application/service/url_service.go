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
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/url/cache"
	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
	"github.com/PraveenGongada/shortly/internal/domain/url/repository"
	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
	"github.com/PraveenGongada/shortly/internal/shared/utils"
)

// URLService defines the interface for URL use cases
type URLService interface {
	CreateShortURL(
		ctx context.Context,
		userID string,
		req *valueobject.CreateURLRequest,
	) (*valueobject.CreateURLResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	GetAnalytics(ctx context.Context, shortCode string, userID string) (int, error)
	GetPaginatedURLs(
		ctx context.Context,
		userID string,
		limit int,
		offset int,
	) ([]valueobject.URLResponse, error)
	UpdateURL(ctx context.Context, urlID string, userID string, newURL string) error
	DeleteURL(ctx context.Context, urlID string, userID string) error
}

type urlService struct {
	generator  interfaces.ShortCodeGenerator
	validator  interfaces.URLValidator
	repository repository.URLRepository
	cache      cache.URLCache
	logger     logger.Logger
	maxRetries int
}

func NewURLService(
	generator interfaces.ShortCodeGenerator,
	validator interfaces.URLValidator,
	repository repository.URLRepository,
	cache cache.URLCache,
	logger logger.Logger,
	maxRetries int,
) URLService {
	if maxRetries <= 0 {
		maxRetries = 1
	}
	return &urlService{
		generator:  generator,
		validator:  validator,
		repository: repository,
		cache:      cache,
		logger:     logger,
		maxRetries: maxRetries,
	}
}

func (s *urlService) CreateShortURL(
	ctx context.Context,
	userID string,
	req *valueobject.CreateURLRequest,
) (*valueobject.CreateURLResponse, error) {
	s.logger.Info(ctx, "Processing create short URL request",
		logger.String("service", "URLService"),
		logger.String("operation", "CreateShortURL"),
		logger.String("userID", userID),
		logger.String("longURL", req.LongURL))

	url, err := s.createShortURLWithRetries(ctx, userID, req, s.maxRetries)
	if err != nil {
		s.logger.Error(ctx, "Failed to create short URL after retries", logger.Error(err))
		return nil, err
	}

	urlResponse := valueobject.CreateShortURLResponse(url)

	s.logger.Info(ctx, "Short URL creation successful",
		logger.String("shortCode", url.ShortCode()),
		logger.String("urlID", url.ID()))

	return &urlResponse, nil
}

func (s *urlService) createShortURLWithRetries(
	ctx context.Context,
	userID string,
	req *valueobject.CreateURLRequest,
	retriesLeft int,
) (*entity.URL, error) {
	if retriesLeft <= 0 {
		return nil, errors.InternalError("max retries exceeded")
	}

	shortCode, err := s.generator.GenerateShortCode()
	if err != nil {
		return nil, errors.InternalError("short code generation failed")
	}

	exists, err := s.repository.ExistsByShortCode(ctx, shortCode)
	if err != nil {
		return nil, errors.InternalError("database query failed")
	}
	if exists {
		// Retry with a new short code
		return s.createShortURLWithRetries(ctx, userID, req, retriesLeft-1)
	}

	// Generate UUID for new URL
	urlID := utils.GenerateRandomUUID()

	// Create new URL entity (validation happens in domain)
	url, err := entity.NewURL(urlID, userID, shortCode, req.LongURL, s.validator)
	if err != nil {
		return nil, errors.ValidationError(err.Error())
	}

	savedURL, err := s.repository.Save(ctx, url)
	if err != nil {
		return nil, errors.InternalError("save operation failed")
	}

	return savedURL, nil
}

func (s *urlService) GetOriginalURL(
	ctx context.Context,
	shortCode string,
) (string, error) {
	s.logger.Info(ctx, "Processing get original URL request",
		logger.String("service", "URLService"),
		logger.String("operation", "GetOriginalURL"),
		logger.String("shortCode", shortCode))

	// Try to get from cache first
	if cachedURL, err := s.cache.GetOriginalURL(ctx, shortCode); err == nil && cachedURL != "" {
		s.logger.Debug(ctx, "URL found in cache",
			logger.String("shortCode", shortCode))
		s.repository.IncrementRedirects(ctx, shortCode)
		return cachedURL, nil
	}

	// Find URL by short code in repository
	url, err := s.repository.FindByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.Warn(ctx, "URL not found",
			logger.String("shortCode", shortCode))
		return "", errors.NotFoundError("URL not found")
	}

	// Cache the result for future requests
	s.cache.SetShortURL(ctx, shortCode, url.LongURL(), time.Minute*5)

	s.repository.IncrementRedirects(ctx, shortCode)

	s.logger.Info(ctx, "URL retrieval successful",
		logger.String("shortCode", shortCode))

	return url.LongURL(), nil
}

func (s *urlService) GetAnalytics(
	ctx context.Context,
	shortCode string,
	userID string,
) (int, error) {
	url, err := s.repository.FindByShortCode(ctx, shortCode)
	if err != nil {
		return -1, errors.NotFoundError("URL not found")
	}

	// Check ownership using domain method
	if !url.IsOwnedBy(userID) {
		return -1, errors.UnauthorizedError("not authorized to view analytics")
	}

	return url.Redirects(), nil
}

func (s *urlService) GetPaginatedURLs(
	ctx context.Context,
	userID string,
	limit int,
	offset int,
) ([]valueobject.URLResponse, error) {
	urls, err := s.repository.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.InternalError("query failed")
	}

	return valueobject.CreateGetURLsResponse(urls), nil
}

func (s *urlService) UpdateURL(
	ctx context.Context,
	urlID string,
	userID string,
	newURL string,
) error {
	url, err := s.repository.FindByID(ctx, urlID)
	if err != nil {
		return errors.NotFoundError("URL not found")
	}

	// Check ownership using domain method
	if !url.IsOwnedBy(userID) {
		return errors.UnauthorizedError("not authorized to update this URL")
	}

	// Update URL using domain method (includes validation)
	if err := url.UpdateLongURL(newURL, s.validator); err != nil {
		return errors.ValidationError(err.Error())
	}

	// Invalidate cache
	s.cache.InvalidateShortURL(ctx, url.ShortCode())

	// Save changes
	return s.repository.Update(ctx, url)
}

func (s *urlService) DeleteURL(
	ctx context.Context,
	urlID string,
	userID string,
) error {
	url, err := s.repository.FindByID(ctx, urlID)
	if err != nil {
		return errors.NotFoundError("URL not found")
	}

	// Check ownership using domain method
	if !url.IsOwnedBy(userID) {
		return errors.UnauthorizedError("not authorized to delete this URL")
	}

	// Invalidate cache
	s.cache.InvalidateShortURL(ctx, url.ShortCode())

	return s.repository.Delete(ctx, urlID, userID)
}
