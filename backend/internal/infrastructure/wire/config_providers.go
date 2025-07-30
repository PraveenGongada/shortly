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
	"github.com/PraveenGongada/shortly/internal/application/service"
	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	urlCache "github.com/PraveenGongada/shortly/internal/domain/url/cache"
	urlRepository "github.com/PraveenGongada/shortly/internal/domain/url/repository"
	urlDomainService "github.com/PraveenGongada/shortly/internal/domain/url/service"
	infraConfig "github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/PraveenGongada/shortly/internal/infrastructure/logging/zerolog"
)

// Config providers that create domain config interfaces
func ProvideDatabaseConfig() config.DatabaseConfig {
	cfg := infraConfig.GetGlobalConfig()
	secrets := infraConfig.GetGlobalSecrets()
	return infraConfig.NewDatabaseConfigAdapter(cfg, secrets)
}

func ProvideAuthConfig() config.AuthConfig {
	cfg := infraConfig.GetGlobalConfig()
	secrets := infraConfig.GetGlobalSecrets()
	return infraConfig.NewAuthConfigAdapter(cfg, secrets)
}

func ProvideURLConfig() config.URLConfig {
	cfg := infraConfig.GetGlobalConfig()
	return infraConfig.NewURLConfigAdapter(cfg)
}

func ProvideServerConfig() config.ServerConfig {
	cfg := infraConfig.GetGlobalConfig()
	return infraConfig.NewServerConfigAdapter(cfg)
}

func ProvideLogConfig() config.LogConfig {
	cfg := infraConfig.GetGlobalConfig()
	return infraConfig.NewLogConfigAdapter(cfg)
}

func ProvideLogger(logConfig config.LogConfig) logger.Logger {
	return zerolog.New()
}

func NewGenerator(urlConfig config.URLConfig) interfaces.ShortCodeGenerator {
	return urlDomainService.NewGenerator(urlConfig.ShortURLLength())
}

func NewURLService(
	generator interfaces.ShortCodeGenerator,
	validator interfaces.URLValidator,
	repository urlRepository.URLRepository,
	cache urlCache.URLCache,
	logger logger.Logger,
	urlConfig config.URLConfig,
) service.URLService {
	return service.NewURLService(generator, validator, repository, cache, logger, urlConfig.MaxCollisionRetries())
}
