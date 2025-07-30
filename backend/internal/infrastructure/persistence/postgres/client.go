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
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

type Store interface {
	Pool() *pgxpool.Pool
	GetQueryTimeout() time.Duration
}

type store struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func (s *store) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *store) GetQueryTimeout() time.Duration {
	return s.queryTimeout
}

func NewPostgresClient(log logger.Logger, dbConfig config.DatabaseConfig) Store {
	log.Info(context.Background(), "Initializing postgres connection...")

	dbHost := dbConfig.Host()
	dbPort := dbConfig.Port()
	dbName := dbConfig.Name()
	dbUser := dbConfig.User()
	dbPassword := dbConfig.Password()
	sslMode := dbConfig.SSLMode()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		sslMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error(context.Background(), "Error parsing database config", logger.Error(err))
		panic(err) // Fatal error during initialization
	}

	poolConfig.MaxConns = dbConfig.MaxConnections()
	poolConfig.MinConns = dbConfig.MinConnections()
	poolConfig.MaxConnLifetime = dbConfig.MaxConnLifetime()
	poolConfig.MaxConnIdleTime = dbConfig.MaxConnIdleTime()
	poolConfig.HealthCheckPeriod = dbConfig.HealthCheckPeriod()

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnectTimeout())
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error(context.Background(), "Error creating connection pool", logger.Error(err))
		panic(err) // Fatal error during initialization
	}

	if err = pool.Ping(ctx); err != nil {
		log.Error(context.Background(), "Error connecting to Database", logger.Error(err))
		panic(err) // Fatal error during initialization
	}

	log.Info(context.Background(), "Success connecting to DB with pgx")

	return &store{
		pool:         pool,
		queryTimeout: dbConfig.QueryTimeout(),
	}
}
