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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/shared/rsa"
)

type JwtToken interface {
	SignRSA(claims jwt.MapClaims) valueobject.Token
}

type JwtTokenImpl struct {
	jwtTokenTimeExp time.Duration
}

func NewJwt() JwtToken {
	jwtTokenDuration, err := time.ParseDuration(config.Get().Auth.JwtToken.Expired)
	if err != nil {
		log.Err(err).Msg(config.Get().Auth.JwtToken.Expired)
	}
	return &JwtTokenImpl{
		jwtTokenTimeExp: jwtTokenDuration,
	}
}

func (o JwtTokenImpl) SignRSA(claims jwt.MapClaims) valueobject.Token {
	timeNow := time.Now()
	tokenExpired := timeNow.Add(o.jwtTokenTimeExp).Unix()

	if claims["id"] == nil {
		return valueobject.Token{}
	}

	token := jwt.New(jwt.SigningMethodRS256)
	// setup userdata
	_, checkExp := claims["exp"]
	_, checkIat := claims["iat"]

	// if user didn't define claims expired
	if !checkExp {
		claims["exp"] = tokenExpired
	}
	// if user didn't define claims iat
	if !checkIat {
		claims["iat"] = timeNow.Unix()
	}
	claims["token_type"] = "access_token"

	token.Claims = claims
	authToken := new(valueobject.Token)
	privateRsa, err := rsa.ReadPrivateKeyFromEnv(config.Get().Application.Key.Rsa.Private)
	if err != nil {
		log.Err(err).Msg("err read private key rsa from env")
		return valueobject.Token{}
	}
	tokenString, err := token.SignedString(privateRsa)
	if err != nil {
		log.Err(err).Msg("err read private rsa")
		return valueobject.Token{}
	}

	authToken.Token = tokenString
	authToken.Type = "Bearer"

	return *authToken
}
