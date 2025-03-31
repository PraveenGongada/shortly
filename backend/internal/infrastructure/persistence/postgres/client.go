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
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresClient() *PostgresStore {
	log.Info().Msg("Initializing postgres connection...")

	dbConfig := config.Get().DB.Postgres
	dbHost := dbConfig.Host
	dbPort := dbConfig.Port
	dbName := dbConfig.Name
	dbUser := dbConfig.User
	dbPassword := dbConfig.Pass
	sslMode := dbConfig.SSLMode

	sHost := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		sslMode,
	)
	db, err := sql.Open("postgres", sHost)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading Database")
	}

	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Database")
	}

	log.Info().Str("Name", dbName).Msg("Success connecting to DB")
	return &PostgresStore{DB: db}
}
