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
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

func RequestLogger(domainLogger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx := r.Context()
			// Add request ID to logger context using domain logger
			ctx = domainLogger.WithContext(ctx, logger.String("request_id", requestID))
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-ID", requestID)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			path := r.URL.Path
			method := r.Method
			statusCode := ww.Status()
			bytesWritten := ww.BytesWritten()
			duration := time.Since(start)
			remoteAddr := r.RemoteAddr

			// Use domain logger with structured fields
			logFields := []logger.Field{
				logger.String("request_id", requestID),
				logger.String("remote_addr", remoteAddr),
				logger.String("method", method),
				logger.String("path", path),
				logger.Int("status", statusCode),
				logger.Int("bytes_written", bytesWritten),
				logger.Any("duration", duration),
			}

			if statusCode >= 500 {
				domainLogger.Error(ctx, "Request completed with server error", logFields...)
			} else if statusCode >= 400 {
				domainLogger.Warn(ctx, "Request completed with client error", logFields...)
			} else {
				domainLogger.Info(ctx, "Request completed", logFields...)
			}
		})
	}
}
