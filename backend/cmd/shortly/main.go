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

package main

import (
	"context"

	_ "github.com/PraveenGongada/shortly/api"
	"github.com/PraveenGongada/shortly/internal/application/service"
	urlDomainService "github.com/PraveenGongada/shortly/internal/domain/url/service"
	userDomainService "github.com/PraveenGongada/shortly/internal/domain/user/service"
	"github.com/PraveenGongada/shortly/internal/infrastructure/auth"
	"github.com/PraveenGongada/shortly/internal/infrastructure/cache/redis"
	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/cookie"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/handler"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
	"github.com/PraveenGongada/shortly/internal/infrastructure/logging/zerolog"
	"github.com/PraveenGongada/shortly/internal/infrastructure/persistence/postgres"
	"github.com/PraveenGongada/shortly/internal/shared/graceful"
)

// @title URL Shortener API
// @version 1.0
// @description A simple URL shortener service with user authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// Load configuration and create adapters
	cfg := config.GetGlobalConfig()
	secrets := config.GetGlobalSecrets()
	logConfig := config.NewLogConfigAdapter(cfg)
	dbConfig := config.NewDatabaseConfigAdapter(cfg, secrets)
	redisConfig := config.NewRedisConfigAdapter(cfg, secrets)
	securityConfig := config.NewSecurityConfigAdapter(cfg)
	authConfig := config.NewAuthConfigAdapter(cfg, secrets)
	urlConfig := config.NewURLConfigAdapter(cfg)
	serverConfig := config.NewServerConfigAdapter(cfg)

	// Initialize logger with config
	logger := zerolog.InitLogger(logConfig)

	// Initialize domain logger adapter
	domainLogger := zerolog.NewWithLogger(logger)

	// Initialize infrastructure layer
	postgresClient := postgres.NewPostgresClient(domainLogger, dbConfig)
	redisClient := redis.NewClient(domainLogger, redisConfig)
	tokenGenerator := auth.NewJwtTokenGenerator(domainLogger, authConfig)
	cookieManager := cookie.NewCookieManager(authConfig)

	// Initialize cache
	urlCache := redis.NewURLCache(redisClient, domainLogger)

	// Initialize domain services
	urlGenerator := urlDomainService.NewGenerator(urlConfig.ShortURLLength())
	urlValidator := urlDomainService.NewValidator()
	userValidator := userDomainService.NewValidator()
	userHasher := userDomainService.NewHasher()

	// Initialize repositories (infrastructure implementations of domain interfaces)
	userRepository := postgres.NewUserRepository(postgresClient, domainLogger)
	urlRepository := postgres.NewURLRepository(postgresClient, domainLogger)

	// Initialize application services (use case implementations)
	userService := service.NewUserService(
		userValidator,
		userHasher,
		userRepository,
		tokenGenerator,
		domainLogger,
	)
	urlService := service.NewURLService(
		urlGenerator,
		urlValidator,
		urlRepository,
		urlCache,
		domainLogger,
		urlConfig.MaxCollisionRetries(),
	)

	// Initialize HTTP layer
	handlers := handler.New(userService, urlService, cookieManager, domainLogger, authConfig)
	routerInstance := router.New(handlers, securityConfig, domainLogger)
	server := http.New(routerInstance, domainLogger, serverConfig)

	// Set up graceful shutdown with proper context handling
	graceful.GracefulShutdown(
		func() error {
			return server.Listen()
		},
		serverConfig.GracefulShutdownTimeout(),
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
			"postgres": func(ctx context.Context) error {
				postgresClient.Close()
				return nil
			},
			"redis": func(ctx context.Context) error {
				return redisClient.Close()
			},
		},
		domainLogger,
	)
}
