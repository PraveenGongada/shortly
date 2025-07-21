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

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/PraveenGongada/shortly/internal/domain/user/repository"
	userService "github.com/PraveenGongada/shortly/internal/domain/user/service"
	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/auth"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
	"github.com/PraveenGongada/shortly/internal/shared/utils"
)

type UserServiceImpl struct {
	jwtAuth        auth.JwtToken
	userRepository repository.UserRepository
}

func NewUserService(jwtAuth auth.JwtToken, repo repository.UserRepository) userService.UserService {
	return &UserServiceImpl{
		jwtAuth:        jwtAuth,
		userRepository: repo,
	}
}

func (s UserServiceImpl) UserLogin(
	ctx context.Context,
	req *valueobject.UserLoginReqest,
) (*valueobject.UserTokenRespBody, error) {
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Email)))[:12]
	logger := log.Ctx(ctx).With().Str("emailHash", emailHash).Str("operation", "UserLogin").Logger()
	logger.Info().Msg("User login attempt")

	user, err := s.userRepository.FindByEmail(req.Email)
	if err != nil {
		logger.Warn().Err(err).Msg("User lookup failed")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.Warn().Msg("Password verification failed")
		return nil, errors.Unauthorized("cannot find user with given email & password")
	}

	userToken := s.jwtAuth.SignRSA(jwt.MapClaims{
		"id": user.ID,
	})

	tokenRes := valueobject.UserTokenRespBody{
		Type:  userToken.Type,
		Token: userToken.Token,
	}

	logger.Info().Str("userId", user.ID).Msg("Login successful")
	return &tokenRes, nil
}

func (s UserServiceImpl) UserLogout(ctx context.Context) error {
	return nil
}

func (s UserServiceImpl) UserRegister(
	ctx context.Context,
	req *valueobject.UserRegisterRequest,
) (*valueobject.UserTokenRespBody, error) {
	uuid := utils.GenerateRandomUUID()

	bcryptPass, err := utils.BcryptString(req.Password)

	if err != nil {
		return nil, err
	} else {
		req.Password = bcryptPass
	}

	user, err := s.userRepository.CreateUser(req, uuid)
	if err != nil {
		return nil, err
	}

	userToken := s.jwtAuth.SignRSA(jwt.MapClaims{
		"id": user.ID,
	})

	tokenRes := valueobject.UserTokenRespBody{
		Type:  userToken.Type,
		Token: userToken.Token,
	}

	return &tokenRes, nil
}
