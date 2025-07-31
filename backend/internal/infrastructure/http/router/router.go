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

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/handler"
	httpmiddleware "github.com/PraveenGongada/shortly/internal/infrastructure/http/middleware"
)

type Router struct {
	handlers       *handler.Handler
	securityConfig config.SecurityConfig
	logger         logger.Logger
}

func New(handlers *handler.Handler, securityConfig config.SecurityConfig, logger logger.Logger) *Router {
	return &Router{
		handlers:       handlers,
		securityConfig: securityConfig,
		logger:         logger,
	}
}

func (h *Router) Router(r *chi.Mux) {
	r.Use(httpmiddleware.RequestTimeout(h.securityConfig, h.logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   h.securityConfig.AllowedOrigins(),
		AllowedMethods:   h.securityConfig.AllowedMethods(),
		AllowedHeaders:   h.securityConfig.AllowedHeaders(),
		AllowCredentials: h.securityConfig.AllowCredentials(),
		MaxAge:           h.securityConfig.MaxAge(),
	}))

	h.handlers.Router(r)

	// Mount Swagger UI handler
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
