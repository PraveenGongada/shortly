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

package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

type Client interface {
	Client() redis.UniversalClient
	Close() error
}

type client struct {
	rdb redis.UniversalClient
}

func (c *client) Client() redis.UniversalClient {
	return c.rdb
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func NewClient(log logger.Logger, redisConfig config.RedisConfig) Client {
	log.Info(context.Background(), "Initializing Redis connection...")

	addrs := redisConfig.Addrs()
	if len(addrs) == 0 {
		addr := fmt.Sprintf("%s:%s", redisConfig.Host(), strconv.Itoa(redisConfig.Port()))
		addrs = []string{addr}
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           addrs,
		Password:        redisConfig.Password(),
		DB:              redisConfig.Database(),
		DialTimeout:     redisConfig.DialTimeout(),
		ReadTimeout:     redisConfig.ReadTimeout(),
		WriteTimeout:    redisConfig.WriteTimeout(),
		PoolSize:        redisConfig.MaxActive(),
		MinIdleConns:    redisConfig.MaxIdle(),
		ConnMaxIdleTime: redisConfig.IdleTimeout(),
		ConnMaxLifetime: redisConfig.MaxConnLifetime(),
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), redisConfig.DialTimeout())
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error(context.Background(), "Error connecting to Redis", logger.Error(err))
		panic(err)
	}

	log.Info(context.Background(), "Successfully connected to Redis")
	return &client{rdb: rdb}
}
