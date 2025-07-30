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
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

// GenerateRandomUUID generates a new UUID string
func GenerateRandomUUID() string {
	return uuid.New().String()
}

// BcryptString hashes a string using bcrypt
// Deprecated: Password hashing should be done in the domain layer
func BcryptString(value string) (string, error) {
	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	if err != nil {
		return "", errors.InternalError("password hashing failed")
	}
	return string(bcryptPass), nil
}
