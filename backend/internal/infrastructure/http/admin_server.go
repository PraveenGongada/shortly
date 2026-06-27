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

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

type AdminServer struct {
	serverConfig config.ServerConfig
	logger       logger.Logger
	httpServer   *http.Server
}

func NewAdminServer(serverConfig config.ServerConfig, logger logger.Logger) *AdminServer {
	return &AdminServer{
		serverConfig: serverConfig,
		logger:       logger,
	}
}

func (a *AdminServer) Listen() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	port := a.serverConfig.AdminPort()

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	a.logger.Info(
		context.Background(),
		"Admin server started",
		logger.String("port", fmt.Sprintf(":%d", port)),
	)
	return a.httpServer.ListenAndServe()
}

func (a *AdminServer) Shutdown(ctx context.Context) error {
	if a.httpServer == nil {
		return nil
	}
	return a.httpServer.Shutdown(ctx)
}
