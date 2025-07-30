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

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
)

type Server struct {
	HttpRouter   *router.Router
	httpServer   *http.Server
	logger       logger.Logger
	serverConfig config.ServerConfig
}

func New(HttpRouter *router.Router, log logger.Logger, serverConfig config.ServerConfig) *Server {
	return &Server{
		HttpRouter:   HttpRouter,
		logger:       log,
		serverConfig: serverConfig,
	}
}

func (r *Server) setupRouter(app *chi.Mux) {
	r.HttpRouter.Router(app)
}

func (p *Server) Listen() error {
	app := chi.NewRouter()

	p.setupRouter(app)

	serverPort := fmt.Sprintf(":%d", p.serverConfig.Port())
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	p.logger.Info(context.Background(), "Server started on port", logger.String("port", serverPort))
	return p.httpServer.ListenAndServe()
}

func (p *Server) Shutdown(ctx context.Context) error {
	if err := p.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
