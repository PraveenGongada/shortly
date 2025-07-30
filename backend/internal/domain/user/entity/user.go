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
	"strings"
	"time"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
)

// User represents a user aggregate root
type User struct {
	id        string
	email     string
	password  string
	name      string
	createdAt time.Time
	updatedAt *time.Time
}

// NewUser creates a new user with validation
func NewUser(id, email, password, name string, validator interfaces.UserValidator, hasher interfaces.PasswordHasher) (*User, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}
	if err := validator.ValidateName(name); err != nil {
		return nil, err
	}

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        id,
		email:     strings.ToLower(strings.TrimSpace(email)),
		password:  hashedPassword,
		name:      strings.TrimSpace(name),
		createdAt: time.Now().UTC(),
	}, nil
}

// NewUserFromRepository creates user from repository data (already validated)
func NewUserFromRepository(id, email, password, name string, createdAt time.Time, updatedAt *time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// VerifyPassword checks if the provided password matches the user's password
func (u *User) VerifyPassword(password string, hasher interfaces.PasswordHasher) error {
	return hasher.VerifyPassword(u.password, password)
}

// UpdateEmail updates the user's email with validation
func (u *User) UpdateEmail(newEmail string, validator interfaces.UserValidator) error {
	if err := validator.ValidateEmail(newEmail); err != nil {
		return err
	}
	u.email = strings.ToLower(strings.TrimSpace(newEmail))
	u.markUpdated()
	return nil
}

// UpdateName updates the user's name with validation
func (u *User) UpdateName(newName string, validator interfaces.UserValidator) error {
	if err := validator.ValidateName(newName); err != nil {
		return err
	}
	u.name = strings.TrimSpace(newName)
	u.markUpdated()
	return nil
}

// ChangePassword changes the user's password with validation
func (u *User) ChangePassword(newPassword string, validator interfaces.UserValidator, hasher interfaces.PasswordHasher) error {
	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.password = hashedPassword
	u.markUpdated()
	return nil
}

// Getters
func (u *User) ID() string            { return u.id }
func (u *User) Email() string         { return u.email }
func (u *User) Name() string          { return u.name }
func (u *User) CreatedAt() time.Time  { return u.createdAt }
func (u *User) UpdatedAt() *time.Time { return u.updatedAt }

// HashedPassword returns the hashed password for persistence
func (u *User) HashedPassword() string { return u.password }

func (u *User) markUpdated() {
	now := time.Now().UTC()
	u.updatedAt = &now
}

type UserList []*User
