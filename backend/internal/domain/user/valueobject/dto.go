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

package valueobject

import "github.com/PraveenGongada/shortly/internal/domain/user/entity"

// Token represents an authentication token
type Token struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// LoginRequest represents user login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents user registration request data
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// TokenResponse represents authentication response
type TokenResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// CreateUserResponse creates a UserResponse from a User entity
func CreateUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}
}
