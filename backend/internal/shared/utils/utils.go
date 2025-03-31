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

package utils

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
)

func GenerateRandomUUID() string {
	uuid := uuid.New().String()
	return uuid
}

func BcryptString(value string) (string, error) {
	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	if err != nil {
		log.Err(err).Msg("Error bcrypting string")
		return "", errors.InternalServerError()
	} else {
		return string(bcryptPass), nil
	}
}

func GenerateShortId() string {
	shortUrlLen := config.Get().Application.ShortUrlLength
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortKey := make([]byte, shortUrlLen)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortKey)
}

func SetCookie(w http.ResponseWriter, token string) error {
	timeNow := time.Now()
	expirationConfig, err := time.ParseDuration(config.Get().Auth.JwtToken.Expired)
	if err == nil {
		expiry := timeNow.Add(expirationConfig)
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expiry,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	} else {
		return err
	}

	return nil
}

func InvalidateCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	}
	http.SetCookie(w, cookie)
}
