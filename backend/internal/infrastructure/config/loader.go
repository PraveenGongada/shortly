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
	Level string `yaml:"level" mapstructure:"LEVEL" validate:"required,oneof=trace debug info warn error fatal panic"`
}

type GracefulConfig struct {
	MaxSecond time.Duration `yaml:"max_second" mapstructure:"MAX_SECOND" validate:"required,min=1s"`
}

type ApplicationConfig struct {
	Port                int            `yaml:"port" mapstructure:"PORT" validate:"required,min=1,max=65535"`
	ShortUrlLength      int8           `yaml:"short_url_length" mapstructure:"SHORT_URL_LENGTH" validate:"required,min=4,max=20"`
	MaxCollisionRetries int8           `yaml:"max_collision_retries" mapstructure:"MAX_COLLISION_RETRIES" validate:"required,min=1,max=10"`
	Environment         string         `yaml:"environment" mapstructure:"ENVIRONMENT" validate:"required,oneof=DEVELOPMENT STAGING PRODUCTION"`
	Log                 LogConfig      `yaml:"log" mapstructure:"LOG"`
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
	SSLMode        string             `yaml:"ssl_mode" mapstructure:"SSL" validate:"required,oneof=disable require verify-ca verify-full"`
	Pool           PostgresPoolConfig `yaml:"pool" mapstructure:"POOL"`
	QueryTimeout   time.Duration      `yaml:"query_timeout" mapstructure:"QUERY_TIMEOUT" validate:"required"`
	ConnectTimeout time.Duration      `yaml:"connect_timeout" mapstructure:"CONNECT_TIMEOUT" validate:"required"`
}

type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"HOST" validate:"required"`
	Port     int    `yaml:"port" mapstructure:"PORT" validate:"required,min=1,max=65535"`
	Database int    `yaml:"database" mapstructure:"DATABASE" validate:"min=0,max=15"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres" mapstructure:"POSTGRES"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"REDIS"`
}

type Config struct {
	Application ApplicationConfig `yaml:"application" mapstructure:"APPLICATION"`
	Auth        AuthConfig        `yaml:"auth" mapstructure:"AUTH"`
	DB          DatabaseConfig    `yaml:"db" mapstructure:"DB"`
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

	l.setDefaults(v)
	l.configureViper(v)

	// Try to read main config file first
	if err := v.ReadInConfig(); err != nil {
		// If main config file not found, try to read modular config files
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := l.loadModularConfigs(v); err != nil {
				// If both main config and modular configs fail, continue with defaults
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
	configFiles := []string{"application", "auth", "database"}
	
	for _, configFile := range configFiles {
		configViper := viper.New()
		configViper.SetConfigName(configFile)
		configViper.SetConfigType("yaml")
		
		for _, path := range l.configPaths {
			configViper.AddConfigPath(path)
		}
		
		if err := configViper.ReadInConfig(); err != nil {
			// Continue if file doesn't exist
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				continue
			}
			return fmt.Errorf("failed to read config file %s: %w", configFile, err)
		}
		
		// Merge settings from this config file
		if err := v.MergeConfigMap(configViper.AllSettings()); err != nil {
			return fmt.Errorf("failed to merge config from %s: %w", configFile, err)
		}
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

func (l *Loader) setDefaults(v *viper.Viper) {
	// Application defaults
	v.SetDefault("APPLICATION.PORT", 8080)
	v.SetDefault("APPLICATION.SHORT_URL_LENGTH", 7)
	v.SetDefault("APPLICATION.MAX_COLLISION_RETRIES", 3)
	v.SetDefault("APPLICATION.ENVIRONMENT", "DEVELOPMENT")
	v.SetDefault("APPLICATION.LOG.LEVEL", "info")
	v.SetDefault("APPLICATION.GRACEFUL.MAX_SECOND", "30s")

	// Auth defaults
	v.SetDefault("AUTH.JWT_TOKEN.TYPE", "Bearer")
	v.SetDefault("AUTH.JWT_TOKEN.EXPIRED", "24h")
	v.SetDefault("AUTH.JWT_TOKEN.REFRESH_EXPIRED", "168h")

	// Database defaults
	v.SetDefault("DB.POSTGRES.HOST", "localhost")
	v.SetDefault("DB.POSTGRES.PORT", 5432)
	v.SetDefault("DB.POSTGRES.NAME", "shortly")
	v.SetDefault("DB.POSTGRES.SSL", "disable")
	v.SetDefault("DB.POSTGRES.QUERY_TIMEOUT", "5s")
	v.SetDefault("DB.POSTGRES.CONNECT_TIMEOUT", "10s")

	// Pool defaults
	v.SetDefault("DB.POSTGRES.POOL.MAX_CONNS", 25)
	v.SetDefault("DB.POSTGRES.POOL.MIN_CONNS", 5)
	v.SetDefault("DB.POSTGRES.POOL.MAX_CONN_LIFETIME", "1h")
	v.SetDefault("DB.POSTGRES.POOL.MAX_CONN_IDLE_TIME", "15m")
	v.SetDefault("DB.POSTGRES.POOL.HEALTH_CHECK_PERIOD", "2m")

	// Redis defaults
	v.SetDefault("DB.REDIS.HOST", "localhost")
	v.SetDefault("DB.REDIS.PORT", 6379)
	v.SetDefault("DB.REDIS.DATABASE", 0)
}