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
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/user/entity"
	"github.com/PraveenGongada/shortly/internal/domain/user/repository"
	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
	"github.com/PraveenGongada/shortly/internal/shared/utils"
)

// TokenGenerator defines interface for token generation
type TokenGenerator interface {
	GenerateToken(userID string) (string, string, error) // returns type, token, error
}

// UserService defines the interface for user use cases
type UserService interface {
	Login(ctx context.Context, req *valueobject.LoginRequest) (*valueobject.TokenResponse, error)
	Logout(ctx context.Context) error
	Register(ctx context.Context, req *valueobject.RegisterRequest) (*valueobject.TokenResponse, error)
}

type userService struct {
	validator      interfaces.UserValidator
	hasher         interfaces.PasswordHasher
	repository     repository.UserRepository
	tokenGenerator TokenGenerator
	logger         logger.Logger
}

func NewUserService(
	validator interfaces.UserValidator,
	hasher interfaces.PasswordHasher,
	repository repository.UserRepository,
	tokenGenerator TokenGenerator,
	logger logger.Logger,
) UserService {
	return &userService{
		validator:      validator,
		hasher:         hasher,
		repository:     repository,
		tokenGenerator: tokenGenerator,
		logger:         logger,
	}
}

func (s *userService) Login(
	ctx context.Context,
	req *valueobject.LoginRequest,
) (*valueobject.TokenResponse, error) {
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Email)))[:12]
	s.logger.Info(ctx, "User login attempt",
		logger.String("emailHash", emailHash),
		logger.String("operation", "Login"))

	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.UnauthorizedError("invalid email or password")
	}

	if err := user.VerifyPassword(req.Password, s.hasher); err != nil {
		s.logger.Warn(ctx, "Password verification failed",
			logger.String("emailHash", emailHash))
		return nil, errors.UnauthorizedError("invalid email or password")
	}

	tokenType, token, err := s.tokenGenerator.GenerateToken(user.ID())
	if err != nil {
		s.logger.Error(ctx, "Token generation failed",
			logger.String("userID", user.ID()),
			logger.Error(err))
		return nil, errors.InternalError("token generation failed")
	}

	tokenRes := &valueobject.TokenResponse{
		Type:  tokenType,
		Token: token,
	}

	s.logger.Info(ctx, "Login successful",
		logger.String("userID", user.ID()))
	return tokenRes, nil
}

func (s *userService) Logout(ctx context.Context) error {
	return nil
}

func (s *userService) Register(
	ctx context.Context,
	req *valueobject.RegisterRequest,
) (*valueobject.TokenResponse, error) {

	// Check if user already exists
	exists, err := s.repository.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.InternalError("database query failed")
	}
	if exists {
		return nil, errors.ConflictError("email already registered")
	}

	// Generate UUID for new user
	userID := utils.GenerateRandomUUID()

	// Create new user entity (validation and password hashing happen in domain)
	user, err := entity.NewUser(userID, req.Email, req.Password, req.Name, s.validator, s.hasher)
	if err != nil {
		s.logger.Warn(ctx, "User validation failed",
			logger.String("operation", "Register"),
			logger.Error(err))
		return nil, errors.ValidationError(err.Error())
	}

	// Save user
	savedUser, err := s.repository.Save(ctx, user)
	if err != nil {
		return nil, errors.InternalError("save operation failed")
	}

	// Generate token
	tokenType, token, err := s.tokenGenerator.GenerateToken(savedUser.ID())
	if err != nil {
		s.logger.Error(ctx, "Token generation failed",
			logger.String("userID", savedUser.ID()),
			logger.Error(err))
		return nil, errors.InternalError("token generation failed")
	}

	s.logger.Info(ctx, "User registered successfully",
		logger.String("userID", savedUser.ID()))

	return &valueobject.TokenResponse{
		Type:  tokenType,
		Token: token,
	}, nil
}
