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
	"github.com/PraveenGongada/shortly/internal/shared/health"
)

type Server struct {
	HttpRouter   *router.Router
	httpServer   *http.Server
	logger       logger.Logger
	serverConfig config.ServerConfig
	readiness    *health.Checker
}

func New(
	HttpRouter *router.Router,
	log logger.Logger,
	serverConfig config.ServerConfig,
	readiness *health.Checker,
) *Server {
	return &Server{
		HttpRouter:   HttpRouter,
		logger:       log,
		serverConfig: serverConfig,
		readiness:    readiness,
	}
}

func (r *Server) setupRouter(app *chi.Mux) {
	r.HttpRouter.Router(app)
}

func (p *Server) withProbes(app http.Handler) http.Handler {
	probes := newProbes(p.logger, p.readiness)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/livez":
			probes.Live(w, r)
		case "/readyz":
			probes.Ready(w, r)
		default:
			app.ServeHTTP(w, r)
		}
	})
}

func (p *Server) Listen() error {
	app := chi.NewRouter()

	p.setupRouter(app)

	serverPort := fmt.Sprintf(":%d", p.serverConfig.Port())
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: p.withProbes(app),
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
