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
	"log"

	_ "github.com/PraveenGongada/shortly/api"
	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
	"github.com/PraveenGongada/shortly/internal/infrastructure/logging/zerolog"
	wireProviders "github.com/PraveenGongada/shortly/internal/infrastructure/wire"
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
	// Load configuration
	cfg := config.GetGlobalConfig()
	logConfig := config.NewLogConfigAdapter(cfg)
	serverConfig := config.NewServerConfigAdapter(cfg)
	securityConfig := config.NewSecurityConfigAdapter(cfg)

	// Initialize logger with config
	logger := zerolog.InitLogger(logConfig)

	// Initialize domain logger adapter
	domainLogger := zerolog.NewWithLogger(logger)

	// Initialize the complete application using Wire
	app, err := wireProviders.InitializeApplication(domainLogger)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Initialize HTTP layer (these remain manual as they're infrastructure wiring)
	routerInstance := router.New(app.Handler, securityConfig, domainLogger)
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
				app.PostgresClient.Close()
				return nil
			},
			"redis": func(ctx context.Context) error {
				return app.RedisClient.Close()
			},
		},
		domainLogger,
	)
}
