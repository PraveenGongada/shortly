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

package logging

import (
	"io"
	"os"
	"time"

	"github.com/PraveenGongada/shortly/internal/infrastructure/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ToZerologLevel(l string) zerolog.Level {
	switch l {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func GetDefaultLogLevel(env string) zerolog.Level {
	switch env {
	case "DEVELOPMENT":
		return zerolog.DebugLevel
	case "PRODUCTION":
		return zerolog.InfoLevel
	case "MASTER":
		return zerolog.WarnLevel
	default:
		return zerolog.InfoLevel
	}
}

func InitLogger() {
	env := config.Get().Application.Environment
	configuredLogLevel := config.Get().Application.Log.Level

	zerolog.TimeFieldFormat = time.RFC3339

	var output io.Writer = os.Stdout

	if env == "DEVELOPMENT" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	var logLevel zerolog.Level
	if configuredLogLevel != "" {
		logLevel = ToZerologLevel(configuredLogLevel)
	} else {
		logLevel = GetDefaultLogLevel(env)
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	log.Info().
		Str("environment", env).
		Str("log_level", logLevel.String()).
		Msg("Zerolog initialized...")
}
