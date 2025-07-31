//go:build wireinject
// +build wireinject

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

package wire

import (
	"github.com/google/wire"

	"github.com/PraveenGongada/shortly/internal/application/service"
	urlDomainService "github.com/PraveenGongada/shortly/internal/domain/url/service"
	userDomainService "github.com/PraveenGongada/shortly/internal/domain/user/service"
	"github.com/PraveenGongada/shortly/internal/infrastructure/auth"
	"github.com/PraveenGongada/shortly/internal/infrastructure/cache/redis"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/cookie"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/handler"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/router"
	"github.com/PraveenGongada/shortly/internal/infrastructure/persistence/postgres"
)

var DomainLayerSet = wire.NewSet(
	NewGenerator,
	urlDomainService.NewValidator,
	userDomainService.NewValidator,
	userDomainService.NewHasher,

	ProvideURLConfig,
)

var ApplicationLayerSet = wire.NewSet(
	service.NewUserService,
	NewURLService,
)

var InterfaceLayerSet = wire.NewSet(
	handler.New,
	router.New,
)

var InfrastructureLayerSet = wire.NewSet(
	ProvideDatabaseConfig,
	ProvideAuthConfig,
	ProvideServerConfig,
	ProvideLogConfig,
	ProvideRedisConfig,
	ProvideSecurityConfig,
	postgres.NewPostgresClient,
	postgres.NewUserRepository,
	postgres.NewURLRepository,
	NewRedisClient,
	redis.NewURLCache,
	auth.NewJwtTokenGenerator,
	wire.Bind(new(service.TokenGenerator), new(auth.TokenGenerator)),

	cookie.NewCookieManager,
)

var FullApplicationSet = wire.NewSet(
	DomainLayerSet,
	InfrastructureLayerSet,
	ApplicationLayerSet,
	InterfaceLayerSet,
)
