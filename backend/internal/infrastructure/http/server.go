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

package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
)

type HttpImpl struct {
	HttpRouter *router.HttpRouterImpl
	httpServer *http.Server
}

func NewHttpProtocol(HttpRouter *router.HttpRouterImpl) *HttpImpl {
	return &HttpImpl{
		HttpRouter: HttpRouter,
	}
}

func (r *HttpImpl) setupRouter(app *chi.Mux) {
	r.HttpRouter.Router(app)
}

func (p *HttpImpl) Listen() error {
	app := chi.NewRouter()

	p.setupRouter(app)

	serverPort := fmt.Sprintf(":%d", config.Get().Application.Port)
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	log.Info().Msgf("Server started on Port %s ", serverPort)
	return p.httpServer.ListenAndServe()
}

func (p *HttpImpl) Shutdown(ctx context.Context) error {
	if err := p.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
