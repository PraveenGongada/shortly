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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

func JwtAuth(log logger.Logger, authConfig config.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug(r.Context(), "Verifying JWT token",
				logger.String("middleware", "JwtVerifyToken"))

			var JwtToken string

			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				JwtToken = strings.TrimPrefix(authHeader, "Bearer ")
				log.Debug(r.Context(), "Token found in Authorization header",
					logger.String("middleware", "JwtVerifyToken"))
			}

			if JwtToken == "" {
				cookie, err := r.Cookie("token")
				if err != nil || cookie.Value == "" {
					log.Warn(r.Context(), "No token found in request",
						logger.String("middleware", "JwtVerifyToken"),
						logger.Error(err))
					response.Json(w, http.StatusUnauthorized, "Token is empty", nil)
					return
				}
				JwtToken = cookie.Value
				log.Debug(r.Context(), "Token found in cookie",
					logger.String("middleware", "JwtVerifyToken"))
			}

			token, err := jwt.Parse(JwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					log.Warn(r.Context(), "Unexpected signing method",
						logger.String("middleware", "JwtVerifyToken"),
						logger.String("method", token.Header["alg"].(string)))
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				publicRsa := authConfig.GetRSAPublicKey()
				if publicRsa == nil {
					log.Error(r.Context(), "RSA public key is not configured",
						logger.String("middleware", "JwtVerifyToken"))
					return nil, errors.InternalError("RSA key configuration error")
				}
				return publicRsa, nil
			})

			if err != nil || !token.Valid {
				log.Warn(r.Context(), "Token is not valid",
					logger.String("middleware", "JwtVerifyToken"),
					logger.Error(err))
				response.Json(w, http.StatusUnauthorized, "Token is not valid", nil)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Warn(r.Context(), "Invalid token claims format",
					logger.String("middleware", "JwtVerifyToken"))
				response.Json(w, http.StatusUnauthorized, "Invalid token claims", nil)
				return
			}

			id, ok := claims["id"].(string)
			if !ok || id == "" {
				log.Warn(r.Context(), "User ID not found in token",
					logger.String("middleware", "JwtVerifyToken"))
				response.Json(w, http.StatusUnauthorized, "ID not found", nil)
				return
			}

			rawExp, ok := claims["exp"].(float64)
			if !ok {
				log.Warn(r.Context(), "Expiration time missing or invalid",
					logger.String("middleware", "JwtVerifyToken"))
				response.Json(w, http.StatusUnauthorized, "Expiration time missing or invalid", nil)
				return
			}
			exp := int64(rawExp)
			if exp < time.Now().Unix() {
				log.Warn(r.Context(), "Token has expired",
					logger.String("middleware", "JwtVerifyToken"),
					logger.Any("expired", time.Unix(exp, 0)),
					logger.Any("now", time.Now()))
				response.Json(w, http.StatusUnauthorized, "Token has expired", nil)
				return
			}

			r.Header.Set("id", id)
			log.Info(r.Context(), "JWT token validated successfully",
				logger.String("middleware", "JwtVerifyToken"),
				logger.String("userId", id))

			next.ServeHTTP(w, r)
		})
	}
}
