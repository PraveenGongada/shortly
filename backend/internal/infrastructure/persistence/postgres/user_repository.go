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

	"github.com/lib/pq"
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
	query := `SELECT * from "user" where email=$1`
	rows, err := r.DB.Query(query, email)
	if err != nil {
		log.Err(err).Msg("error fetching user")
		return nil, errors.InternalServerError()
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, errors.NotFound(fmt.Sprintf("cannot find user with email %v", email))
	}

	user := &entity.User{}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Err(err).Msg("Error scanning rows in FindByEmail")
			return nil, errors.InternalServerError()
		}
	}

	if user.ID == "" {
		return nil, errors.Unauthorized("cannot find user with given email & password")
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindByID(ID string) (*entity.User, error) {
	query := `SELECT * from "user" where id=$1`

	rows, err := r.DB.Query(query, ID)
	if err != nil {
		log.Err(err).Msg("error fetching user")
		return nil, errors.InternalServerError()
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, errors.NotFound(fmt.Sprintf("cannot find user with email %v", ID))
	}

	user := &entity.User{}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Err(err).Msg("Error scanning rows in FindById")
			return nil, errors.InternalServerError()
		}
	}

	return user, nil
}

func (r UserRepositoryImpl) CreateUser(
	req *valueobject.UserRegisterRequest,
	uuid string,
) (*entity.User, error) {
	query := `INSERT INTO "user" (id, name, email, password) VALUES ($1,$2,$3,$4) RETURNING *;`

	rows, err := r.DB.Query(query, uuid, req.Name, req.Email, req.Password)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == "unique_violation" {
			return nil, errors.Unauthorized("email is already registered")
		} else {
			log.Err(err).Msg("Error creating user")
			return nil, errors.InternalServerError()
		}
	}

	defer rows.Close()

	user := &entity.User{}

	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Err(err).Msg("Error scanning rows CreateUser")
			return nil, errors.InternalServerError()
		}
	}

	return user, nil
}
