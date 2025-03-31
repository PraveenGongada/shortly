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
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
	"github.com/PraveenGongada/shortly/internal/shared/rsa"
)

func JwtVerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var JwtToken string

		// Extract Bearer token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			JwtToken = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// If Authorization header is missing, check cookie
		if JwtToken == "" {
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				response.Json(w, http.StatusUnauthorized, "Token is empty", nil)
				return
			}
			JwtToken = cookie.Value
		}

		// Parse JWT token
		token, err := jwt.Parse(JwtToken, func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is RSA
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Load RSA public key
			publicRsa, err := rsa.ReadPublicKeyFromEnv(config.Get().Application.Key.Rsa.Public)
			if err != nil {
				log.Err(err).Msg("Error reading RSA public key from env")
				return nil, errors.InternalServerError()
			}
			return publicRsa, nil
		})

		// Validate token
		if err != nil || !token.Valid {
			log.Err(err).Msg("Token is not valid")
			response.Json(w, http.StatusUnauthorized, "Token is not valid", nil)
			return
		}

		// Extract claims safely
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Json(w, http.StatusUnauthorized, "Invalid token claims", nil)
			return
		}

		// Extract user ID
		id, ok := claims["id"].(string)
		if !ok || id == "" {
			response.Json(w, http.StatusUnauthorized, "ID not found", nil)
			return
		}

		// Check expiration time
		rawExp, ok := claims["exp"].(float64)
		if !ok {
			response.Json(w, http.StatusUnauthorized, "Expiration time missing or invalid", nil)
			return
		}
		exp := int64(rawExp)
		if exp < time.Now().Unix() {
			response.Json(w, http.StatusUnauthorized, "Token has expired", nil)
			return
		}

		// Set user ID in request header for next handler
		r.Header.Set("id", id)

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
