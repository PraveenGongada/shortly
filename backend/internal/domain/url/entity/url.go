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

package entity

import (
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
)

// URL represents a URL aggregate root
type URL struct {
	id        string
	userID    string
	shortCode string
	longURL   string
	redirects int
	createdAt time.Time
	updatedAt *time.Time
}

// NewURL creates a new URL with validation
func NewURL(id, userID, shortCode, longURL string, validator interfaces.URLValidator) (*URL, error) {
	if err := validator.ValidateUserID(userID); err != nil {
		return nil, err
	}
	if err := validator.ValidateShortCode(shortCode); err != nil {
		return nil, err
	}
	if err := validator.ValidateURL(longURL); err != nil {
		return nil, err
	}

	return &URL{
		id:        id,
		userID:    userID,
		shortCode: shortCode,
		longURL:   longURL,
		redirects: 0,
		createdAt: time.Now().UTC(),
	}, nil
}

// NewURLFromRepository creates URL from repository data (already validated)
func NewURLFromRepository(id, userID, shortCode, longURL string, redirects int, createdAt time.Time, updatedAt *time.Time) *URL {
	return &URL{
		id:        id,
		userID:    userID,
		shortCode: shortCode,
		longURL:   longURL,
		redirects: redirects,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// UpdateLongURL updates the target URL with validation
func (u *URL) UpdateLongURL(newURL string, validator interfaces.URLValidator) error {
	if err := validator.ValidateURL(newURL); err != nil {
		return err
	}
	u.longURL = newURL
	u.markUpdated()
	return nil
}

// IncrementRedirects increases the redirect count
func (u *URL) IncrementRedirects() {
	u.redirects++
	u.markUpdated()
}

// IsOwnedBy checks if the URL belongs to the specified user
func (u *URL) IsOwnedBy(userID string) bool {
	return u.userID == userID
}

// Getters
func (u *URL) ID() string            { return u.id }
func (u *URL) UserID() string        { return u.userID }
func (u *URL) ShortCode() string     { return u.shortCode }
func (u *URL) LongURL() string       { return u.longURL }
func (u *URL) Redirects() int        { return u.redirects }
func (u *URL) CreatedAt() time.Time  { return u.createdAt }
func (u *URL) UpdatedAt() *time.Time { return u.updatedAt }

func (u *URL) markUpdated() {
	now := time.Now().UTC()
	u.updatedAt = &now
}

type URLList []*URL
