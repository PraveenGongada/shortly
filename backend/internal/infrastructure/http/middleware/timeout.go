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

package httpmiddleware

import (
	"context"
	"net/http"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
)

func RequestTimeout(securityConfig config.SecurityConfig, log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeout := securityConfig.RequestTimeout()
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})
			panicChan := make(chan interface{}, 1)

			go func() {
				defer func() {
					if p := recover(); p != nil {
						panicChan <- p
					}
				}()

				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case p := <-panicChan:
				panic(p)
			case <-done:
				return
			case <-ctx.Done():
				log.Warn(ctx, "Request timeout exceeded",
					logger.String("middleware", "RequestTimeout"),
					logger.String("method", r.Method),
					logger.String("path", r.URL.Path),
					logger.Any("timeout", timeout))

				if ctx.Err() == context.DeadlineExceeded {
					response.Json(w, http.StatusGatewayTimeout, "Request timeout", nil)
					return
				}
			}
		})
	}
}
