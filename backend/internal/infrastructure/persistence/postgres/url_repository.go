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
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
	"github.com/PraveenGongada/shortly/internal/domain/url/repository"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

type urlRepository struct {
	store  Store
	logger logger.Logger
}

// NewURLRepository creates a new URL repository implementation
func NewURLRepository(store Store, logger logger.Logger) repository.URLRepository {
	return &urlRepository{
		store:  store,
		logger: logger,
	}
}

func (r *urlRepository) Save(ctx context.Context, url *entity.URL) (*entity.URL, error) {

	query := `INSERT INTO "url" (id, user_id, short_url, long_url, redirects, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, user_id, short_url, long_url, redirects, created_at, updated_at`

	var id, userID, shortCode, longURL string
	var redirects int
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query,
		url.ID(),
		url.UserID(),
		url.ShortCode(),
		url.LongURL(),
		url.Redirects(),
		url.CreatedAt(),
	).Scan(&id, &userID, &shortCode, &longURL, &redirects, &createdAt, &updatedAt)

	if err != nil {
		r.logger.Error(ctx, "Error saving URL", 
			logger.String("urlId", url.ID()),
			logger.String("operation", "Save"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	savedURL := entity.NewURLFromRepository(id, userID, shortCode, longURL, redirects, createdAt, updatedAt)
	r.logger.Info(ctx, "URL saved successfully",
		logger.String("urlId", url.ID()),
		logger.String("operation", "Save"))
	return savedURL, nil
}

func (r *urlRepository) FindByShortCode(ctx context.Context, shortCode string) (*entity.URL, error) {

	query := `SELECT id, user_id, short_url, long_url, redirects, created_at, updated_at 
			  FROM "url" WHERE short_url = $1`

	var id, userID, scannedShortCode, longURL string
	var redirects int
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query, shortCode).Scan(
		&id, &userID, &scannedShortCode, &longURL, &redirects, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Debug(ctx, "URL not found",
				logger.String("shortCode", shortCode),
				logger.String("operation", "FindByShortCode"))
			return nil, errors.NotFoundError("URL not found")
		}
		r.logger.Error(ctx, "Error finding URL by short code",
			logger.String("shortCode", shortCode),
			logger.String("operation", "FindByShortCode"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	url := entity.NewURLFromRepository(id, userID, scannedShortCode, longURL, redirects, createdAt, updatedAt)
	r.logger.Debug(ctx, "URL found successfully",
		logger.String("shortCode", shortCode),
		logger.String("operation", "FindByShortCode"))
	return url, nil
}

func (r *urlRepository) FindByID(ctx context.Context, id string) (*entity.URL, error) {

	query := `SELECT id, user_id, short_url, long_url, redirects, created_at, updated_at 
			  FROM "url" WHERE id = $1`

	var urlID, userID, shortCode, longURL string
	var redirects int
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query, id).Scan(
		&urlID, &userID, &shortCode, &longURL, &redirects, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Debug(ctx, "URL not found",
				logger.String("urlId", id),
				logger.String("operation", "FindByID"))
			return nil, errors.NotFoundError("URL not found")
		}
		r.logger.Error(ctx, "Error finding URL by ID",
			logger.String("urlId", id),
			logger.String("operation", "FindByID"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	url := entity.NewURLFromRepository(urlID, userID, shortCode, longURL, redirects, createdAt, updatedAt)
	r.logger.Debug(ctx, "URL found successfully",
		logger.String("urlId", id),
		logger.String("operation", "FindByID"))
	return url, nil
}

func (r *urlRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.URL, error) {

	query := `SELECT id, user_id, short_url, long_url, redirects, created_at, updated_at 
			  FROM "url" 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC 
			  LIMIT $2 OFFSET $3`

	rows, err := r.store.Pool().Query(ctx, query, userID, limit, offset)
	if err != nil {
		r.logger.Error(ctx, "Error querying URLs by user ID",
			logger.String("userId", userID),
			logger.Int("limit", limit),
			logger.Int("offset", offset),
			logger.String("operation", "FindByUserID"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}
	defer rows.Close()

	var urls []*entity.URL

	for rows.Next() {
		var id, uid, shortCode, longURL string
		var redirects int
		var createdAt time.Time
		var updatedAt *time.Time

		err := rows.Scan(&id, &uid, &shortCode, &longURL, &redirects, &createdAt, &updatedAt)
		if err != nil {
			r.logger.Error(ctx, "Error scanning URL row",
				logger.String("userId", userID),
				logger.String("operation", "FindByUserID"),
				logger.Error(err))
			return nil, errors.InternalError("database operation failed")
		}

		url := entity.NewURLFromRepository(id, uid, shortCode, longURL, redirects, createdAt, updatedAt)
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, "Error iterating URL rows",
			logger.String("userId", userID),
			logger.String("operation", "FindByUserID"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	r.logger.Debug(ctx, "URLs found successfully",
		logger.String("userId", userID),
		logger.Int("count", len(urls)),
		logger.String("operation", "FindByUserID"))
	return urls, nil
}

func (r *urlRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM "url" WHERE short_url = $1)`

	var exists bool
	err := r.store.Pool().QueryRow(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		r.logger.Error(ctx, "Error checking if short code exists",
			logger.String("shortCode", shortCode),
			logger.String("operation", "ExistsByShortCode"),
			logger.Error(err))
		return false, errors.InternalError("database operation failed")
	}

	return exists, nil
}

func (r *urlRepository) Update(ctx context.Context, url *entity.URL) error {

	query := `UPDATE "url" 
			  SET long_url = $1, updated_at = $2 
			  WHERE id = $3`

	cmdTag, err := r.store.Pool().Exec(ctx, query,
		url.LongURL(),
		time.Now().UTC(),
		url.ID(),
	)

	if err != nil {
		r.logger.Error(ctx, "Error updating URL",
			logger.String("urlId", url.ID()),
			logger.String("operation", "Update"),
			logger.Error(err))
		return errors.InternalError("database operation failed")
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Debug(ctx, "URL not found for update",
			logger.String("urlId", url.ID()),
			logger.String("operation", "Update"))
		return errors.NotFoundError("URL not found")
	}

	r.logger.Info(ctx, "URL updated successfully",
		logger.String("urlId", url.ID()),
		logger.String("operation", "Update"))
	return nil
}

func (r *urlRepository) Delete(ctx context.Context, id, userID string) error {

	query := `DELETE FROM "url" WHERE id = $1 AND user_id = $2`

	cmdTag, err := r.store.Pool().Exec(ctx, query, id, userID)
	if err != nil {
		r.logger.Error(ctx, "Error deleting URL",
			logger.String("urlId", id),
			logger.String("userId", userID),
			logger.String("operation", "Delete"),
			logger.Error(err))
		return errors.InternalError("database operation failed")
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Debug(ctx, "URL not found for deletion or unauthorized",
			logger.String("urlId", id),
			logger.String("userId", userID),
			logger.String("operation", "Delete"))
		return errors.NotFoundError("URL not found or not authorized")
	}

	r.logger.Info(ctx, "URL deleted successfully",
		logger.String("urlId", id),
		logger.String("userId", userID),
		logger.String("operation", "Delete"))
	return nil
}

func (r *urlRepository) IncrementRedirects(ctx context.Context, shortCode string) error {

	query := `UPDATE "url" 
			  SET redirects = redirects + 1, updated_at = $1 
			  WHERE short_url = $2`

	cmdTag, err := r.store.Pool().Exec(ctx, query, time.Now().UTC(), shortCode)
	if err != nil {
		r.logger.Error(ctx, "Error incrementing redirects",
			logger.String("shortCode", shortCode),
			logger.String("operation", "IncrementRedirects"),
			logger.Error(err))
		return errors.InternalError("database operation failed")
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Debug(ctx, "URL not found for redirect increment",
			logger.String("shortCode", shortCode),
			logger.String("operation", "IncrementRedirects"))
		return errors.NotFoundError("URL not found")
	}

	r.logger.Debug(ctx, "Redirects incremented successfully",
		logger.String("shortCode", shortCode),
		logger.String("operation", "IncrementRedirects"))
	return nil
}
