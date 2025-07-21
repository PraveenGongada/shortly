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
	"crypto/sha256"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/user/entity"
	"github.com/PraveenGongada/shortly/internal/domain/user/repository"
	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
)

type UserRepositoryImpl struct {
	*PostgresStore
}

func NewUserRepository(db *PostgresStore) repository.UserRepository {
	return &UserRepositoryImpl{
		PostgresStore: db,
	}
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	// Hash email for privacy-safe logging - same approach as in service layer
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))[:12]
	logger := log.With().Str("emailHash", emailHash).Str("operation", "FindByEmail").Logger()
	logger.Debug().Msg("Finding user by email")

	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `SELECT id, name, email, password, created_at, updated_at FROM "user" WHERE email=$1`

	user := &entity.User{}
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug().Msg("User not found")
			return nil, errors.Unauthorized("cannot find user with given email & password")
		}
		logger.Error().Err(err).Msg("Error finding user by email")
		return nil, errors.InternalServerError()
	}

	logger.Debug().Str("userId", user.ID).Msg("User found successfully")
	return user, nil
}

func (r *UserRepositoryImpl) FindByID(ID string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `SELECT id, name, email, password, created_at, updated_at FROM "user" WHERE id=$1`

	user := &entity.User{}
	err := r.DB.QueryRow(ctx, query, ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("cannot find user with ID %v", ID))
		}
		log.Err(err).Msg("Error finding user by ID")
		return nil, errors.InternalServerError()
	}

	return user, nil
}

func (r UserRepositoryImpl) CreateUser(
	req *valueobject.UserRegisterRequest,
	uuid string,
) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.GetQueryTimeout())
	defer cancel()

	query := `INSERT INTO "user" (id, name, email, password) VALUES ($1,$2,$3,$4) RETURNING id, name, email, password, created_at, updated_at`

	user := &entity.User{}
	err := r.DB.QueryRow(ctx, query, uuid, req.Name, req.Email, req.Password).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" { // unique_violation
			return nil, errors.Unauthorized("email is already registered")
		}
		log.Err(err).Msg("Error creating user")
		return nil, errors.InternalServerError()
	}

	return user, nil
}
