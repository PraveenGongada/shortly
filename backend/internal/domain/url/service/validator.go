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

package service

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
)

const (
	MaxURLLength       = 2048
	MinShortCodeLength = 4
	MaxShortCodeLength = 12
)

type validator struct{}

func NewValidator() interfaces.URLValidator {
	return &validator{}
}

func (v *validator) ValidateURL(longURL string) error {
	longURL = strings.TrimSpace(longURL)
	if longURL == "" {
		return errors.New("URL cannot be empty")
	}
	if len(longURL) > MaxURLLength {
		return errors.New("URL cannot exceed 2048 characters")
	}

	// Parse URL to validate format
	parsedURL, err := url.Parse(longURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Ensure it has a scheme
	if parsedURL.Scheme == "" {
		return errors.New("URL must include a scheme (http:// or https://)")
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must use http or https scheme")
	}

	// Ensure it has a host
	if parsedURL.Host == "" {
		return errors.New("URL must include a valid host")
	}

	return nil
}

func (v *validator) ValidateShortCode(shortCode string) error {
	shortCode = strings.TrimSpace(shortCode)
	if shortCode == "" {
		return errors.New("short code cannot be empty")
	}
	if len(shortCode) < MinShortCodeLength || len(shortCode) > MaxShortCodeLength {
		return errors.New("short code must be between 4 and 12 characters")
	}

	// Check for valid characters (alphanumeric only)
	for _, char := range shortCode {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return errors.New("short code can only contain alphanumeric characters")
		}
	}

	return nil
}

func (v *validator) ValidateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("user ID cannot be empty")
	}
	return nil
}
