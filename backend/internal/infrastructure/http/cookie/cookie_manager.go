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

package cookie

import (
	"net/http"
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
)

// Manager defines cookie management operations
type Manager interface {
	SetAuthCookie(w http.ResponseWriter, token string) error
	InvalidateAuthCookie(w http.ResponseWriter)
}

// cookieManager implements cookie management
type cookieManager struct {
	domain     string
	secure     bool
	httpOnly   bool
	sameSite   http.SameSite
	expiration time.Duration
}

// NewCookieManager creates a new cookie manager
func NewCookieManager(authConfig config.AuthConfig) Manager {
	expiration, err := time.ParseDuration(authConfig.JWTTokenExpiry())
	if err != nil {
		expiration = 24 * time.Hour // default to 24 hours
	}

	return &cookieManager{
		domain:     "", // Set from config if needed
		secure:     true,
		httpOnly:   true,
		sameSite:   http.SameSiteNoneMode,
		expiration: expiration,
	}
}

func (cm *cookieManager) SetAuthCookie(w http.ResponseWriter, token string) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(cm.expiration),
		HttpOnly: cm.httpOnly,
		Secure:   cm.secure,
		SameSite: cm.sameSite,
		Path:     "/",
		Domain:   cm.domain,
	}
	http.SetCookie(w, cookie)
	return nil
}

func (cm *cookieManager) InvalidateAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Domain:  cm.domain,
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	}
	http.SetCookie(w, cookie)
}
