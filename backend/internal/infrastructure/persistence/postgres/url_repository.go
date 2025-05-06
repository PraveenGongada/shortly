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
	"fmt"

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
	query := `INSERT INTO "url" VALUES ($1,$2,$3,$4) RETURNING *;`

	rows, err := r.DB.Query(query, id, userId, shortUrl, req.LongUrl)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	url := &entity.Url{}

	for rows.Next() {
		err = rows.Scan(
			&url.Id,
			&url.UserID,
			&url.ShortUrl,
			&url.LongUrl,
			&url.Redirects,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			log.Err(err).Msg("Error scanning rows in CreateShortUrl")
			return nil, errors.InternalServerError()
		}
	}
	return url, nil
}

func (r UrlRepositoryImpl) GetLongUrl(shortUrl string) (string, error) {
	logger := log.With().Str("shortUrl", shortUrl).Str("operation", "GetLongUrl").Logger()
	logger.Debug().Msg("Updating redirect count and retrieving long URL")

	updateAnalytics := `UPDATE url SET redirects=redirects+1 where short_url=$1`

	_, err := r.DB.Exec(updateAnalytics, shortUrl)
	if err != nil {
		log.Err(err).Msg("Error in updating the redirect count")
		return "", errors.InternalServerError()
	}

	longUrl := ""

	query := `SELECT (long_url) FROM "url" WHERE short_url = $1;`

	row := r.DB.QueryRow(query, shortUrl)

	if row.Err() != nil {
		log.Err(row.Err()).Msg("error in GetLongUrl Repo")
		return "", errors.InternalServerError()
	}

	row.Scan(
		&longUrl,
	)
	logger.Debug().Str("longUrl", longUrl).Msg("Successfully retrieved long URL")

	return longUrl, nil
}

func (r UrlRepositoryImpl) GetAnalytics(shortUrl string, userId string) (int, error) {
	query := `SELECT user_id, redirects FROM url where short_url = $1;`

	row, err := r.DB.Query(query, shortUrl)
	if err != nil {
		log.Err(err).Msg("Error in GetAnalytics -> Repository")
		return -1, errors.InternalServerError()
	}

	var user string
	var redirects int

	if row.Next() {
		err = row.Scan(&user, &redirects)
		if err != nil {
			log.Err(err).Msg("Error scanning row in GetAnalytics -> Repository")
			return -1, errors.InternalServerError()
		}
	} else {
		log.Info().Msg("No rows found for the given short URL")
		return -1, errors.InternalServerError()
	}

	if user != userId {
		return -1, errors.Unauthorized("unauthorized request")
	}

	return redirects, nil
}

func (r *UrlRepositoryImpl) GetPaginatedUrls(
	userId string,
	limit int,
	offset int,
) ([]entity.Url, error) {
	query := `
        SELECT id, short_url, long_url, redirects 
        FROM url 
        WHERE user_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := r.DB.Query(query, userId, limit, offset)
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
			log.Err(err).Str("layer", "repository").Msg("Error GetPaginatedUrls")
			return nil, errors.InternalServerError()
		}
		url.Redirects = redirects
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		log.Err(err).Str("layer", "repository").Msg("Error GetPaginatedUrls")
		return nil, errors.InternalServerError()
	}

	return urls, nil
}

func (r UrlRepositoryImpl) UpdateUrl(urlId string, newUrl string) error {
	query := `UPDATE url SET long_url=$1 WHERE id=$2;`

	_, err := r.DB.Exec(query, newUrl, urlId)
	if err != nil {
		log.Err(err).Msg("Error in UpdateUrl -> Repository")
		return errors.InternalServerError()
	}

	return nil
}

func (r UrlRepositoryImpl) DeleteUrl(urlId string, userId string) error {
	query := `DELETE FROM url WHERE id=$1 AND user_id=$2;`

	res, err := r.DB.Exec(query, urlId, userId)
	if err != nil {
		log.Err(err).Msg("Error in DeleteUrl  -> Repository")
		return errors.InternalServerError()
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		log.Err(fmt.Errorf("Error in DeleteUrl -> Repository"))
		return errors.Unauthorized("Cannot delete the url")
	}

	return nil
}
