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

package config

import (
	"crypto/rsa"
	"time"

	domainConfig "github.com/PraveenGongada/shortly/internal/domain/shared/config"
)

type DatabaseConfigAdapter struct {
	config  *Config
	secrets SecretProvider
}

func NewDatabaseConfigAdapter(cfg *Config, secrets SecretProvider) domainConfig.DatabaseConfig {
	return &DatabaseConfigAdapter{config: cfg, secrets: secrets}
}

func (d *DatabaseConfigAdapter) Host() string          { return d.config.DB.Postgres.Host }
func (d *DatabaseConfigAdapter) Port() int             { return d.config.DB.Postgres.Port }
func (d *DatabaseConfigAdapter) Name() string          { return d.config.DB.Postgres.Name }
func (d *DatabaseConfigAdapter) User() string          { return d.secrets.GetDatabaseUser() }
func (d *DatabaseConfigAdapter) Password() string      { return d.secrets.GetDatabasePassword() }
func (d *DatabaseConfigAdapter) SSLMode() string       { return d.config.DB.Postgres.SSLMode }
func (d *DatabaseConfigAdapter) MaxConnections() int32 { return d.config.DB.Postgres.Pool.MaxConns }
func (d *DatabaseConfigAdapter) MinConnections() int32 { return d.config.DB.Postgres.Pool.MinConns }
func (d *DatabaseConfigAdapter) MaxConnLifetime() time.Duration {
	return d.config.DB.Postgres.Pool.MaxConnLifetime
}
func (d *DatabaseConfigAdapter) MaxConnIdleTime() time.Duration {
	return d.config.DB.Postgres.Pool.MaxConnIdleTime
}
func (d *DatabaseConfigAdapter) HealthCheckPeriod() time.Duration {
	return d.config.DB.Postgres.Pool.HealthCheckPeriod
}
func (d *DatabaseConfigAdapter) QueryTimeout() time.Duration {
	return d.config.DB.Postgres.QueryTimeout
}
func (d *DatabaseConfigAdapter) ConnectTimeout() time.Duration {
	return d.config.DB.Postgres.ConnectTimeout
}

type AuthConfigAdapter struct {
	config  *Config
	secrets SecretProvider
}

func NewAuthConfigAdapter(cfg *Config, secrets SecretProvider) domainConfig.AuthConfig {
	return &AuthConfigAdapter{config: cfg, secrets: secrets}
}

func (a *AuthConfigAdapter) JWTTokenType() string              { return a.config.Auth.JwtToken.Type }
func (a *AuthConfigAdapter) JWTTokenExpiry() string            { return a.config.Auth.JwtToken.Expired }
func (a *AuthConfigAdapter) JWTRefreshExpiry() string          { return a.config.Auth.JwtToken.RefreshExpired }
func (a *AuthConfigAdapter) GetRSAPublicKey() *rsa.PublicKey   { return a.secrets.GetRSAPublicKey() }
func (a *AuthConfigAdapter) GetRSAPrivateKey() *rsa.PrivateKey { return a.secrets.GetRSAPrivateKey() }

type URLConfigAdapter struct {
	config *Config
}

func NewURLConfigAdapter(cfg *Config) domainConfig.URLConfig {
	return &URLConfigAdapter{config: cfg}
}

func (u *URLConfigAdapter) ShortURLLength() int { return int(u.config.Application.ShortUrlLength) }
func (u *URLConfigAdapter) MaxCollisionRetries() int {
	return int(u.config.Application.MaxCollisionRetries)
}

type ServerConfigAdapter struct {
	config *Config
}

func NewServerConfigAdapter(cfg *Config) domainConfig.ServerConfig {
	return &ServerConfigAdapter{config: cfg}
}

func (s *ServerConfigAdapter) Port() int { return s.config.Application.Port }
func (s *ServerConfigAdapter) GracefulShutdownTimeout() time.Duration {
	return s.config.Application.Graceful.MaxSecond
}

type LogConfigAdapter struct {
	config *Config
}

func NewLogConfigAdapter(cfg *Config) domainConfig.LogConfig {
	return &LogConfigAdapter{config: cfg}
}

func (l *LogConfigAdapter) Environment() string { return l.config.Application.Environment }
func (l *LogConfigAdapter) Level() string       { return l.config.Application.Log.Level }
