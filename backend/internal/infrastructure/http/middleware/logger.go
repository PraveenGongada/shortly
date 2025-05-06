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
	"github.com/rs/zerolog/log"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := r.Context()
		ctx = log.With().Str("request_id", requestID).Logger().WithContext(ctx)
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

		logger := log.With().
			Str("request_id", requestID).
			Str("remote_addr", remoteAddr).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Int("bytes_written", bytesWritten).
			Dur("duration", duration).
			Logger()

		if statusCode >= 500 {
			logger.Error().Msg("Request completed with server error")
		} else if statusCode >= 400 {
			logger.Warn().Msg("Request completed with client error")
		} else {
			logger.Info().Msg("Request completed")
		}
	})
}
