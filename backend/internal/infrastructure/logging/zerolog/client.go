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

package zerolog

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/PraveenGongada/shortly/internal/domain/shared/config"
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

func InitLogger(logConfig config.LogConfig) zerolog.Logger {
	env := logConfig.Environment()
	configuredLogLevel := logConfig.Level()
	format := logConfig.Format()
	outputPath := logConfig.Output()
	caller := logConfig.Caller()
	timestamp := logConfig.Timestamp()
	timestampFormat := logConfig.TimestampFormat()

	if timestampFormat != "" {
		zerolog.TimeFieldFormat = timestampFormat
	} else {
		zerolog.TimeFieldFormat = time.RFC3339
	}

	if caller {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			return short + ":" + strconv.Itoa(line)
		}
	}

	var output io.Writer = os.Stdout
	if outputPath == "stderr" {
		output = os.Stderr
	} else if outputPath != "stdout" && outputPath != "" {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			output = os.Stdout
		} else {
			output = file
		}
	}

	useConsoleFormat := format == "console"
	if useConsoleFormat {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: timestampFormat,
		}
	}

	var logLevel zerolog.Level
	if configuredLogLevel != "" {
		logLevel = ToZerologLevel(configuredLogLevel)
	} else {
		logLevel = GetDefaultLogLevel(env)
	}
	zerolog.SetGlobalLevel(logLevel)

	loggerCtx := zerolog.New(output).With()
	if timestamp {
		loggerCtx = loggerCtx.Timestamp()
	}
	if caller {
		loggerCtx = loggerCtx.Caller()
	}
	logger := loggerCtx.Logger()

	logger.Info().
		Str("environment", env).
		Str("log_level", logLevel.String()).
		Str("format", format).
		Str("output", outputPath).
		Bool("caller", caller).
		Bool("timestamp", timestamp).
		Msg("Zerolog initialized...")

	return logger
}
