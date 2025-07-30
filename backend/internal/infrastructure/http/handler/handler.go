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

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/PraveenGongada/shortly/internal/application/service"
	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/cookie"
	httpmiddleware "github.com/PraveenGongada/shortly/internal/infrastructure/http/middleware"
)

type Handler struct {
	userService   service.UserService
	urlService    service.URLService
	cookieManager cookie.Manager
	logger        logger.Logger
	authConfig    config.AuthConfig
}

func New(
	userService service.UserService,
	urlService service.URLService,
	cookieManager cookie.Manager,
	logger logger.Logger,
	authConfig config.AuthConfig,
) *Handler {
	return &Handler{
		userService:   userService,
		urlService:    urlService,
		cookieManager: cookieManager,
		logger:        logger,
		authConfig:    authConfig,
	}
}

func (h *Handler) Router(r chi.Router) {
	r.Use(middleware.StripSlashes)
	r.Use(httpmiddleware.RequestLogger(h.logger))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", h.UserLogin)
			r.Post("/register", h.UserRegister)
			r.Get("/logout", h.UserLogout)
		})
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.JwtAuth(h.logger, h.authConfig))
			r.Get("/urls", h.GetPaginatedURLs)
			r.Route("/url", func(r chi.Router) {
				r.Post("/create", h.CreateShortURL)
				r.Patch("/update", h.UpdateURL)
				r.Delete("/{urlId}", h.DeleteURL)
				r.Get("/analytics/{shortUrl}", h.GetAnalytics)
			})
		})
		r.Get("/{shortUrl}", h.GetLongURL)
	})

	r.Group(func(r chi.Router) {
		r.Get("/{shortUrl}", h.RedirectUser)
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("App is running!"))
}
