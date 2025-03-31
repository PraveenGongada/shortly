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

	_ "github.com/PraveenGongada/shortly/docs"
	"github.com/PraveenGongada/shortly/internal/application/service"
	"github.com/PraveenGongada/shortly/internal/infrastructure/auth"
	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/handler"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
	"github.com/PraveenGongada/shortly/internal/infrastructure/logging"
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
	logging.InitLogger()

	postgresClient := postgres.NewPostgresClient()
	jwtAuth := auth.NewJwt()

	// Repositories
	userRepository := postgres.NewUserRepository(postgresClient)
	urlRepository := postgres.NewUrlRepository(postgresClient)

	// Services
	userService := service.NewUserService(jwtAuth, userRepository)
	urlService := service.NewUrlService(urlRepository)

	// HTTP Layer
	handlers := handler.NewHttpHandler(userService, urlService)
	router := router.NewHttpRoute(handlers)
	server := http.NewHttpProtocol(router)

	// Set up graceful shutdown with the server operation
	graceful.GracefulShutdown(
		func() error {
			return server.Listen()
		},
		config.Get().Application.Graceful.MaxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		},
	)
}
