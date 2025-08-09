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
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Level           string `yaml:"level" mapstructure:"LEVEL" validate:"required,oneof=trace debug info warn error fatal panic"`
	Format          string `yaml:"format" mapstructure:"FORMAT" validate:"required,oneof=json console"`
	Output          string `yaml:"output" mapstructure:"OUTPUT" validate:"required"`
	Caller          bool   `yaml:"caller" mapstructure:"CALLER"`
	Timestamp       bool   `yaml:"timestamp" mapstructure:"TIMESTAMP"`
	TimestampFormat string `yaml:"timestamp_format" mapstructure:"TIMESTAMP_FORMAT"`
}

type GracefulConfig struct {
	MaxSecond time.Duration `yaml:"max_second" mapstructure:"MAX_SECOND" validate:"required,min=1s"`
}

type ApplicationConfig struct {
	Port                int            `yaml:"port" mapstructure:"PORT" validate:"required,min=1,max=65535"`
	ShortUrlLength      int8           `yaml:"short_url_length" mapstructure:"SHORT_URL_LENGTH" validate:"required,min=4,max=20"`
	MaxCollisionRetries int8           `yaml:"max_collision_retries" mapstructure:"MAX_COLLISION_RETRIES" validate:"required,min=1,max=10"`
	Environment         string         `yaml:"environment" mapstructure:"ENVIRONMENT" validate:"required,oneof=DEVELOPMENT STAGING PRODUCTION"`
	Graceful            GracefulConfig `yaml:"graceful" mapstructure:"GRACEFUL"`
}

type JwtTokenConfig struct {
	Type           string `yaml:"type" mapstructure:"TYPE" validate:"required"`
	Expired        string `yaml:"expired" mapstructure:"EXPIRED" validate:"required"`
	RefreshExpired string `yaml:"refresh_expired" mapstructure:"REFRESH_EXPIRED" validate:"required"`
}

type AuthConfig struct {
	JwtToken JwtTokenConfig `yaml:"jwt_token" mapstructure:"JWT_TOKEN"`
}

type PostgresPoolConfig struct {
	MaxConns          int32         `yaml:"max_conns" mapstructure:"MAX_CONNS" validate:"required,min=1"`
	MinConns          int32         `yaml:"min_conns" mapstructure:"MIN_CONNS" validate:"required,min=1"`
	MaxConnLifetime   time.Duration `yaml:"max_conn_lifetime" mapstructure:"MAX_CONN_LIFETIME" validate:"required"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" mapstructure:"MAX_CONN_IDLE_TIME" validate:"required"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" mapstructure:"HEALTH_CHECK_PERIOD" validate:"required"`
}

type PostgresConfig struct {
	Host           string             `yaml:"host" mapstructure:"HOST" validate:"required"`
	Port           int                `yaml:"port" mapstructure:"PORT" validate:"required,min=1,max=65535"`
	Name           string             `yaml:"name" mapstructure:"NAME" validate:"required"`
	SSLMode        string             `yaml:"ssl_mode" mapstructure:"SSL_MODE" validate:"required,oneof=disable require verify-ca verify-full"`
	Pool           PostgresPoolConfig `yaml:"pool" mapstructure:"POOL"`
	QueryTimeout   time.Duration      `yaml:"query_timeout" mapstructure:"QUERY_TIMEOUT" validate:"required"`
	ConnectTimeout time.Duration      `yaml:"connect_timeout" mapstructure:"CONNECT_TIMEOUT" validate:"required"`
}

type RedisPoolConfig struct {
	MaxIdle         int           `yaml:"max_idle" mapstructure:"MAX_IDLE" validate:"required,min=1"`
	MaxActive       int           `yaml:"max_active" mapstructure:"MAX_ACTIVE" validate:"required,min=1"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"IDLE_TIMEOUT" validate:"required"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" mapstructure:"MAX_CONN_LIFETIME" validate:"required"`
}

type RedisConfig struct {
	Host         string          `yaml:"host" mapstructure:"HOST" validate:"required"`
	Port         int             `yaml:"port" mapstructure:"PORT" validate:"required,min=1,max=65535"`
	Addrs        []string        `yaml:"addrs" mapstructure:"ADDRS"`
	Database     int             `yaml:"database" mapstructure:"DATABASE" validate:"min=0,max=15"`
	Password     string          `yaml:"password" mapstructure:"PASSWORD"`
	DialTimeout  time.Duration   `yaml:"dial_timeout" mapstructure:"DIAL_TIMEOUT" validate:"required"`
	ReadTimeout  time.Duration   `yaml:"read_timeout" mapstructure:"READ_TIMEOUT" validate:"required"`
	WriteTimeout time.Duration   `yaml:"write_timeout" mapstructure:"WRITE_TIMEOUT" validate:"required"`
	Pool         RedisPoolConfig `yaml:"pool" mapstructure:"POOL"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres" mapstructure:"POSTGRES"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"REDIS"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins" mapstructure:"ALLOWED_ORIGINS" validate:"required"`
	AllowedMethods   []string `yaml:"allowed_methods" mapstructure:"ALLOWED_METHODS" validate:"required"`
	AllowedHeaders   []string `yaml:"allowed_headers" mapstructure:"ALLOWED_HEADERS" validate:"required"`
	AllowCredentials bool     `yaml:"allow_credentials" mapstructure:"ALLOW_CREDENTIALS"`
	MaxAge           int      `yaml:"max_age" mapstructure:"MAX_AGE" validate:"min=0"`
}

type SecurityConfig struct {
	CORS           CORSConfig    `yaml:"cors" mapstructure:"CORS"`
	RequestTimeout time.Duration `yaml:"request_timeout" mapstructure:"REQUEST_TIMEOUT" validate:"required"`
}

type Config struct {
	Application ApplicationConfig `yaml:"application" mapstructure:"APPLICATION"`
	Auth        AuthConfig        `yaml:"auth" mapstructure:"AUTH"`
	Database    DatabaseConfig    `yaml:"database" mapstructure:"DATABASE"`
	Logging     LogConfig         `yaml:"logging" mapstructure:"LOGGING"`
	Security    SecurityConfig    `yaml:"security" mapstructure:"SECURITY"`
}

type Loader struct {
	configPaths []string
	envPrefix   string
}

func NewLoader(configPaths []string, envPrefix string) *Loader {
	if len(configPaths) == 0 {
		configPaths = []string{".", "/etc/shortly/", "configs/"}
	}
	return &Loader{
		configPaths: configPaths,
		envPrefix:   envPrefix,
	}
}

func (l *Loader) Load() (*Config, error) {
	v := viper.New()

	l.configureViper(v)

	// Try to read main config file first
	if err := v.ReadInConfig(); err != nil {
		// If main config file not found, try to read modular config files
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := l.loadModularConfigs(v); err != nil {
				return nil, fmt.Errorf("no configuration files found in paths %v: main config and modular configs both failed: %w", l.configPaths, err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := l.validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func (l *Loader) loadModularConfigs(v *viper.Viper) error {
	configFiles := []string{"application", "auth", "database", "logging", "security"}
	foundFiles := make([]string, 0)
	missingFiles := make([]string, 0)

	for _, configFile := range configFiles {
		configViper := viper.New()
		configViper.SetConfigName(configFile)
		configViper.SetConfigType("yaml")

		for _, path := range l.configPaths {
			configViper.AddConfigPath(path)
		}

		if err := configViper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				missingFiles = append(missingFiles, configFile+".yaml")
				continue
			}
			return fmt.Errorf("failed to read config file %s: %w", configFile, err)
		}

		foundFiles = append(foundFiles, configFile+".yaml")
		// Merge settings from this config file
		if err := v.MergeConfigMap(configViper.AllSettings()); err != nil {
			return fmt.Errorf("failed to merge config from %s: %w", configFile, err)
		}
	}

	if len(missingFiles) > 0 {
		return fmt.Errorf("required configuration files not found: %v (found: %v) in paths: %v", missingFiles, foundFiles, l.configPaths)
	}

	return nil
}

func (l *Loader) configureViper(v *viper.Viper) {
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	for _, path := range l.configPaths {
		v.AddConfigPath(path)
	}

	v.AutomaticEnv()
	if l.envPrefix != "" {
		v.SetEnvPrefix(l.envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func (l *Loader) validate(config *Config) error {
	validate := validator.New()
	return validate.Struct(config)
}
