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

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
	"github.com/PraveenGongada/shortly/internal/domain/url/repository"
	urlService "github.com/PraveenGongada/shortly/internal/domain/url/service"
	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
	"github.com/PraveenGongada/shortly/internal/shared/utils"
)

type UrlServiceImpl struct {
	urlRepository repository.UrlRepository
}

func NewUrlService(repo repository.UrlRepository) urlService.UrlService {
	return &UrlServiceImpl{
		urlRepository: repo,
	}
}

func (s UrlServiceImpl) CreateShortUrl(
	ctx context.Context,
	userId string,
	req *valueobject.CreateUrlRequest,
) (*valueobject.CreateUrlResponse, error) {
	logger := log.Ctx(ctx).With().Str("userId", userId).Logger()
	logger.Debug().Str("longUrl", req.LongUrl).Msg("Starting short URL creation")

	maxCollisionRetries := config.Get().Application.MaxCollisionRetries

	url, err := s.createShotUrlWithRetries(userId, maxCollisionRetries, req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create short URL after retries")
		return nil, err
	}

	urlResponse := valueobject.CreateShortUrlResponse(url)

	logger.Info().
		Str("shortUrl", url.ShortUrl).
		Str("urlId", url.Id).
		Msg("Short URL created successfully")

	return &urlResponse, nil
}

func (s UrlServiceImpl) createShotUrlWithRetries(
	userId string,
	retriesLeft int8,
	req *valueobject.CreateUrlRequest,
) (*entity.Url, error) {
	if retriesLeft <= 0 {
		return nil, errors.InternalServerError()
	}

	shortId := utils.GenerateShortId()
	uuid := utils.GenerateRandomUUID()

	shortUrl, err := s.urlRepository.CreateShortUrl(uuid, userId, shortId, req)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" { // unique_violation
			return s.createShotUrlWithRetries(userId, retriesLeft-1, req)
		} else {
			return nil, errors.InternalServerError()
		}
	}
	return shortUrl, nil
}

func (s UrlServiceImpl) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	logger := log.Ctx(ctx).With().Str("shortUrl", shortUrl).Logger()

	longUrl, err := s.urlRepository.GetLongUrl(shortUrl)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to retrieve long URL")
		return "", err
	}

	logger.Debug().Str("longUrl", longUrl).Msg("Successfully retrieved long URL")
	return longUrl, nil
}

func (s UrlServiceImpl) GetAnalytics(
	ctx context.Context,
	shortUrl string,
	userId string,
) (int, error) {
	count, err := s.urlRepository.GetAnalytics(shortUrl, userId)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (s UrlServiceImpl) GetPaginatedUrls(
	ctx context.Context,
	userId string,
	limit int,
	offset int,
) ([]valueobject.UrlResponse, error) {
	urls, err := s.urlRepository.GetPaginatedUrls(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	urlResponse := valueobject.CreateGetUrlsResponse(urls)

	return urlResponse, nil
}

func (s UrlServiceImpl) UpdateUrl(ctx context.Context, urlId string, newUrl string) error {
	err := s.urlRepository.UpdateUrl(urlId, newUrl)
	if err != nil {
		return err
	}

	return nil
}

func (s UrlServiceImpl) DeleteUrl(ctx context.Context, urlId string, userId string) error {
	err := s.urlRepository.DeleteUrl(urlId, userId)
	if err != nil {
		return err
	}

	return nil
}
