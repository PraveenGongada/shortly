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
)

// DatabaseConfig defines configuration needed for database operations
type DatabaseConfig interface {
	Host() string
	Port() int
	Name() string
	User() string
	Password() string
	SSLMode() string
	MaxConnections() int32
	MinConnections() int32
	MaxConnLifetime() time.Duration
	MaxConnIdleTime() time.Duration
	HealthCheckPeriod() time.Duration
	QueryTimeout() time.Duration
	ConnectTimeout() time.Duration
}

// AuthConfig defines configuration needed for authentication
type AuthConfig interface {
	JWTTokenType() string
	JWTTokenExpiry() string
	JWTRefreshExpiry() string
	GetRSAPublicKey() *rsa.PublicKey
	GetRSAPrivateKey() *rsa.PrivateKey
}

// URLConfig defines configuration needed for URL service
type URLConfig interface {
	ShortURLLength() int
	MaxCollisionRetries() int
}

// ServerConfig defines configuration needed for HTTP server
type ServerConfig interface {
	Port() int
	GracefulShutdownTimeout() time.Duration
}

// LogConfig defines configuration needed for logging
type LogConfig interface {
	Environment() string
	Level() string
	Format() string
	Output() string
	Caller() bool
	Timestamp() bool
	TimestampFormat() string
}

// RedisConfig defines configuration needed for Redis cache
type RedisConfig interface {
	Host() string
	Port() int
	Addrs() []string
	Database() int
	Password() string
	DialTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	MaxIdle() int
	MaxActive() int
	IdleTimeout() time.Duration
	MaxConnLifetime() time.Duration
}

// SecurityConfig defines configuration needed for security features
type SecurityConfig interface {
	AllowedOrigins() []string
	AllowedMethods() []string
	AllowedHeaders() []string
	AllowCredentials() bool
	MaxAge() int
	RequestTimeout() time.Duration
}
