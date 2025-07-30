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
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/user/entity"
	"github.com/PraveenGongada/shortly/internal/domain/user/repository"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

type userRepository struct {
	store  Store
	logger logger.Logger
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(store Store, logger logger.Logger) repository.UserRepository {
	return &userRepository{
		store:  store,
		logger: logger,
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Hash email for privacy-safe logging - same approach as in service layer
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))[:12]
	r.logger.Debug(ctx, "Finding user by email",
		logger.String("emailHash", emailHash),
		logger.String("operation", "FindByEmail"))

	query := `SELECT id, name, email, password, created_at, updated_at FROM "user" WHERE email=$1`

	var id, name, userEmail, password string
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query, email).Scan(
		&id, &name, &userEmail, &password, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Debug(ctx, "User not found",
				logger.String("emailHash", emailHash),
				logger.String("operation", "FindByEmail"))
			return nil, errors.NotFoundError("user not found")
		}
		r.logger.Error(ctx, "Error finding user by email",
			logger.String("emailHash", emailHash),
			logger.String("operation", "FindByEmail"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	user := entity.NewUserFromRepository(id, userEmail, password, name, createdAt, updatedAt)
	r.logger.Debug(ctx, "User found successfully",
		logger.String("emailHash", emailHash),
		logger.String("userId", user.ID()),
		logger.String("operation", "FindByEmail"))
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {

	query := `SELECT id, name, email, password, created_at, updated_at FROM "user" WHERE id=$1`

	var userId, name, email, password string
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query, id).Scan(
		&userId, &name, &email, &password, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Debug(ctx, "User not found",
				logger.String("userId", id),
				logger.String("operation", "FindByID"))
			return nil, errors.NotFoundError("user not found")
		}
		r.logger.Error(ctx, "Error finding user by ID",
			logger.String("userId", id),
			logger.String("operation", "FindByID"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	user := entity.NewUserFromRepository(userId, email, password, name, createdAt, updatedAt)
	r.logger.Debug(ctx, "User found successfully",
		logger.String("userId", id),
		logger.String("operation", "FindByID"))
	return user, nil
}

func (r *userRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {

	query := `INSERT INTO "user" (id, name, email, password, created_at) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, name, email, password, created_at, updated_at`

	var id, name, email, password string
	var createdAt time.Time
	var updatedAt *time.Time

	err := r.store.Pool().QueryRow(ctx, query,
		user.ID(),
		user.Name(),
		user.Email(),
		user.HashedPassword(),
		user.CreatedAt(),
	).Scan(&id, &name, &email, &password, &createdAt, &updatedAt)

	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" { // unique_violation
			r.logger.Warn(ctx, "Email already exists",
				logger.String("userId", user.ID()),
				logger.String("operation", "Save"))
			return nil, errors.ConflictError("email already registered")
		}
		r.logger.Error(ctx, "Error creating user",
			logger.String("userId", user.ID()),
			logger.String("operation", "Save"),
			logger.Error(err))
		return nil, errors.InternalError("database operation failed")
	}

	savedUser := entity.NewUserFromRepository(id, email, password, name, createdAt, updatedAt)
	r.logger.Info(ctx, "User saved successfully",
		logger.String("userId", user.ID()),
		logger.String("operation", "Save"))
	return savedUser, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`

	var exists bool
	err := r.store.Pool().QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.Error(ctx, "Error checking if user exists by email",
			logger.String("operation", "ExistsByEmail"),
			logger.Error(err))
		return false, errors.InternalError("database operation failed")
	}

	return exists, nil
}
