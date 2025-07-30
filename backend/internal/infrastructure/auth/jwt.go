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

package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

// TokenGenerator defines interface for JWT token generation
type TokenGenerator interface {
	GenerateToken(userID string) (string, string, error) // returns type, token, error
}

// JwtTokenGenerator implements token generation using RSA
type JwtTokenGenerator struct {
	tokenExpiration time.Duration
	logger          logger.Logger
	authConfig      config.AuthConfig
}

// NewJwtTokenGenerator creates a new JWT token generator
func NewJwtTokenGenerator(log logger.Logger, authConfig config.AuthConfig) TokenGenerator {
	jwtTokenDuration, err := time.ParseDuration(authConfig.JWTTokenExpiry())
	if err != nil {
		log.Error(context.Background(), "Invalid JWT token expiration config", logger.Error(err))
		jwtTokenDuration = 24 * time.Hour // default to 24 hours
	}
	return &JwtTokenGenerator{
		tokenExpiration: jwtTokenDuration,
		authConfig:      authConfig,
		logger:          log,
	}
}

func (g *JwtTokenGenerator) GenerateToken(userID string) (string, string, error) {
	if userID == "" {
		return "", "", errors.New("user ID cannot be empty")
	}

	timeNow := time.Now()
	tokenExpired := timeNow.Add(g.tokenExpiration).Unix()

	token := jwt.New(jwt.SigningMethodRS256)
	claims := jwt.MapClaims{
		"id":         userID,
		"exp":        tokenExpired,
		"iat":        timeNow.Unix(),
		"token_type": "access_token",
	}

	token.Claims = claims

	privateRsa := g.authConfig.GetRSAPrivateKey()
	if privateRsa == nil {
		g.logger.Error(context.Background(), "Private RSA key not available", 
			logger.String("userID", userID))
		return "", "", errors.New("private RSA key not configured")
	}

	tokenString, err := token.SignedString(privateRsa)
	if err != nil {
		g.logger.Error(context.Background(), "Error signing JWT token",
			logger.String("userID", userID),
			logger.Error(err))
		return "", "", err
	}

	return g.authConfig.JWTTokenType(), tokenString, nil
}
