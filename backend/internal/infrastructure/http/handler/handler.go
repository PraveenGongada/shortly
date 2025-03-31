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

	urlService "github.com/PraveenGongada/shortly/internal/domain/url/service"
	userService "github.com/PraveenGongada/shortly/internal/domain/user/service"
	httpmiddleware "github.com/PraveenGongada/shortly/internal/infrastructure/http/middleware"
)

type HttpHandlerImpl struct {
	userService userService.UserService
	urlService  urlService.UrlService
}

func NewHttpHandler(
	userService userService.UserService,
	urlService urlService.UrlService,
) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		userService: userService,
		urlService:  urlService,
	}
}

func (h *HttpHandlerImpl) Router(r chi.Router) {
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", heartBeatHandler)
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", h.UserLogin)
			r.Post("/register", h.UserRegsiter)
			r.Get("/logout", h.UserLogout)
		})
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.JwtVerifyToken)
			r.Get("/urls", h.GetPaginatedUrls)
			r.Route("/url", func(r chi.Router) {
				r.Post("/create", h.CreateShortUrl)
				r.Patch("/update", h.UpdateUrl)
				r.Delete("/{urlId}", h.DeleteUrl)
				r.Get("/analytics/{shortUrl}", h.GetAnalytics)
			})
		})
		r.Get("/{shortUrl}", h.GetLongUrl)
	})

	r.Group(func(r chi.Router) {
		r.Get("/{shortUrl}", h.RedirectUser)
	})
}

func heartBeatHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("App is running!"))
	return
}
