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

package interfaces

// UserValidator defines the interface for user validation
type UserValidator interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidateName(name string) error
}

// PasswordHasher defines the interface for password hashing operations
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
