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

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
	"github.com/PraveenGongada/shortly/internal/domain/url/repository"
	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
)

type UrlRepositoryImpl struct {
	*PostgresStore
}

func NewUrlRepository(db *PostgresStore) repository.UrlRepository {
	return &UrlRepositoryImpl{
		PostgresStore: db,
	}
}

func (r UrlRepositoryImpl) CreateShortUrl(
	id string,
	userId string,
	shortUrl string,
	req *valueobject.CreateUrlRequest,
) (*entity.Url, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `INSERT INTO "url" (id, user_id, short_url, long_url) VALUES ($1,$2,$3,$4) RETURNING id, user_id, short_url, long_url, redirects, created_at, updated_at`

	url := &entity.Url{}
	err := r.DB.QueryRow(ctx, query, id, userId, shortUrl, req.LongUrl).Scan(
		&url.Id,
		&url.UserID,
		&url.ShortUrl,
		&url.LongUrl,
		&url.Redirects,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		return nil, errors.InternalServerError()
	}
	return url, nil
}

func (r UrlRepositoryImpl) GetLongUrl(shortUrl string) (string, error) {
	logger := log.With().Str("shortUrl", shortUrl).Str("operation", "GetLongUrl").Logger()
	logger.Debug().Msg("Updating redirect count and retrieving long URL")

	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	// Start a transaction to ensure atomicity
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error starting transaction")
		return "", errors.InternalServerError()
	}
	defer tx.Rollback(ctx)

	// Update redirect count
	updateAnalytics := `UPDATE url SET redirects=redirects+1 WHERE short_url=$1`
	_, err = tx.Exec(ctx, updateAnalytics, shortUrl)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating redirect count")
		return "", errors.InternalServerError()
	}

	// Get long URL
	var longUrl string
	query := `SELECT long_url FROM "url" WHERE short_url = $1`
	err = tx.QueryRow(ctx, query, shortUrl).Scan(&longUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Warn().Msg("Short URL not found")
			return "", errors.NotFound("URL not found")
		}
		logger.Error().Err(err).Msg("Error retrieving long URL")
		return "", errors.InternalServerError()
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error().Err(err).Msg("Error committing transaction")
		return "", errors.InternalServerError()
	}

	logger.Debug().Str("longUrl", longUrl).Msg("Successfully retrieved long URL")
	return longUrl, nil
}

func (r UrlRepositoryImpl) GetAnalytics(shortUrl string, userId string) (int, error) {
	logger := log.With().Str("shortUrl", shortUrl).Str("userId", userId).Str("operation", "GetAnalytics").Logger()

	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `SELECT user_id, redirects FROM url WHERE short_url = $1`

	var user string
	var redirects int
	err := r.DB.QueryRow(ctx, query, shortUrl).Scan(&user, &redirects)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug().Msg("URL not found for analytics")
			return -1, errors.NotFound("URL not found")
		}
		logger.Error().Err(err).Msg("Error retrieving analytics data")
		return -1, errors.InternalServerError()
	}

	if user != userId {
		logger.Warn().Msg("Unauthorized analytics access attempt")
		return -1, errors.Unauthorized("unauthorized request")
	}

	return redirects, nil
}

func (r *UrlRepositoryImpl) GetPaginatedUrls(
	userId string,
	limit int,
	offset int,
) ([]entity.Url, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `
        SELECT id, short_url, long_url, redirects 
        FROM url 
        WHERE user_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := r.DB.Query(ctx, query, userId, limit, offset)
	if err != nil {
		log.Err(err).Msg("Error in GetPaginatedUrls -> Repository")
		return nil, errors.InternalServerError()
	}
	defer rows.Close()

	urls := make([]entity.Url, 0, limit)

	for rows.Next() {
		var url entity.Url
		var redirects int

		err := rows.Scan(&url.Id, &url.ShortUrl, &url.LongUrl, &redirects)
		if err != nil {
			log.Err(err).Str("layer", "repository").Msg("Error scanning row in GetPaginatedUrls")
			return nil, errors.InternalServerError()
		}
		url.Redirects = redirects
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		log.Err(err).Str("layer", "repository").Msg("Error with rows in GetPaginatedUrls")
		return nil, errors.InternalServerError()
	}

	return urls, nil
}

func (r UrlRepositoryImpl) UpdateUrl(urlId string, newUrl string) error {
	logger := log.With().Str("urlId", urlId).Str("operation", "UpdateUrl").Logger()

	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `UPDATE url SET long_url=$1 WHERE id=$2`

	cmdTag, err := r.DB.Exec(ctx, query, newUrl, urlId)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating URL")
		return errors.InternalServerError()
	}

	if cmdTag.RowsAffected() == 0 {
		logger.Debug().Msg("URL not found for update")
		return errors.NotFound("URL not found")
	}

	return nil
}

func (r UrlRepositoryImpl) DeleteUrl(urlId string, userId string) error {
	logger := log.With().Str("urlId", urlId).Str("userId", userId).Str("operation", "DeleteUrl").Logger()

	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `DELETE FROM url WHERE id=$1 AND user_id=$2`

	cmdTag, err := r.DB.Exec(ctx, query, urlId, userId)
	if err != nil {
		logger.Error().Err(err).Msg("Error deleting URL")
		return errors.InternalServerError()
	}

	if cmdTag.RowsAffected() == 0 {
		logger.Debug().Msg("URL not found for deletion or unauthorized")
		return errors.Unauthorized("Cannot delete the url")
	}

	return nil
}
