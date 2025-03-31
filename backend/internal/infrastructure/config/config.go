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
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	cfg    Config
	doOnce sync.Once
)

type Config struct {
	Application struct {
		Port                int    `mapstructure:"PORT"`
		ShortUrlLength      int8   `mapstructure:"SHORT_URL_LENGTH"`
		MaxCollisionRetries int8   `mapstructure:"MAX_COLLISSION_RETRIES"`
		Environment         string `mapstructure:"ENVIRONMENT"`
		Log                 struct {
			Path string `mapstructure:"PATH"`
		}
		Key struct {
			Rsa struct {
				Public  string `mapstructure:"PUBLIC"`
				Private string `mapstructure:"PRIVATE"`
			}
		}
		Graceful struct {
			MaxSecond time.Duration `mapstructure:"MAX_SECOND"`
		} `mapstructure:"GRACEFUL"`
	} `mapstructure:"APPLICATION"`

	Auth struct {
		JwtToken struct {
			Type           string `mapstructure:"TYPE"`
			Expired        string `mapstructure:"EXPIRED"`
			RefreshExpired string `mapstructure:"REFRESH_EXPIRED"`
		} `mapstructure:"JWT_TOKEN"`
	} `mapstructure:"AUTH"`

	DB struct {
		Postgres struct {
			Host    string `mapstructure:"HOST"`
			Port    int    `mapstructure:"PORT"`
			Name    string `mapstructure:"NAME"`
			User    string `mapstructure:"USER"`
			Pass    string `mapstructure:"PASS"`
			SSLMode string `mapstructure:"SSL"`
		} `mapstructure:"POSTGRES"`
	} `mapstructure:"DB"`
}

func Get() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/infrastructure/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot read config file")
	}

	doOnce.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("error unmarshaling config")
		}
	})

	return cfg
}
