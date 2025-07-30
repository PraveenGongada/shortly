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
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
)

type adapter struct {
	baseLogger zerolog.Logger
}

func New() logger.Logger {
	return &adapter{
		baseLogger: log.Logger,
	}
}

func NewWithLogger(baseLogger zerolog.Logger) logger.Logger {
	return &adapter{
		baseLogger: baseLogger,
	}
}

func (z *adapter) getLoggerWithContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return z.baseLogger
	}

	// Try to get the logger from context (set by middleware)
	ctxLogger := zerolog.Ctx(ctx)
	if ctxLogger.GetLevel() == zerolog.Disabled {
		return z.baseLogger
	}

	return *ctxLogger
}

func (z *adapter) addFieldToEvent(event *zerolog.Event, key string, value interface{}) *zerolog.Event {
	// Handle special error case first since it has different API (no key parameter)
	if err, ok := value.(error); ok {
		return event.Err(err)
	}

	// Use type-specific methods for common types for better performance
	switch v := value.(type) {
	case string:
		return event.Str(key, v)
	case int:
		return event.Int(key, v)
	case int64:
		return event.Int64(key, v)
	case float64:
		return event.Float64(key, v)
	case bool:
		return event.Bool(key, v)
	case time.Time:
		return event.Time(key, v)
	case time.Duration:
		return event.Dur(key, v)
	case []byte:
		return event.Bytes(key, v)
	default:
		return event.Interface(key, v)
	}
}

func (z *adapter) addFieldToContext(ctx zerolog.Context, key string, value interface{}) zerolog.Context {
	// Handle special error case
	if err, ok := value.(error); ok {
		return ctx.AnErr(key, err)
	}

	// Use type-specific methods for common types for better performance
	switch v := value.(type) {
	case string:
		return ctx.Str(key, v)
	case int:
		return ctx.Int(key, v)
	case int64:
		return ctx.Int64(key, v)
	case float64:
		return ctx.Float64(key, v)
	case bool:
		return ctx.Bool(key, v)
	case time.Time:
		return ctx.Time(key, v)
	case time.Duration:
		return ctx.Dur(key, v)
	case []byte:
		return ctx.Bytes(key, v)
	default:
		return ctx.Interface(key, v)
	}
}

func (z *adapter) addFields(event *zerolog.Event, fields []logger.Field) *zerolog.Event {
	for _, field := range fields {
		event = z.addFieldToEvent(event, field.Key(), field.Value())
	}
	return event
}

func (z *adapter) Debug(ctx context.Context, msg string, fields ...logger.Field) {
	logger := z.getLoggerWithContext(ctx)
	event := logger.Debug().CallerSkipFrame(1)
	event = z.addFields(event, fields)
	event.Msg(msg)
}

func (z *adapter) Info(ctx context.Context, msg string, fields ...logger.Field) {
	logger := z.getLoggerWithContext(ctx)
	event := logger.Info().CallerSkipFrame(1)
	event = z.addFields(event, fields)
	event.Msg(msg)
}

func (z *adapter) Warn(ctx context.Context, msg string, fields ...logger.Field) {
	logger := z.getLoggerWithContext(ctx)
	event := logger.Warn().CallerSkipFrame(1)
	event = z.addFields(event, fields)
	event.Msg(msg)
}

func (z *adapter) Error(ctx context.Context, msg string, fields ...logger.Field) {
	logger := z.getLoggerWithContext(ctx)
	event := logger.Error().CallerSkipFrame(1)
	event = z.addFields(event, fields)
	event.Msg(msg)
}

func (z *adapter) WithContext(ctx context.Context, fields ...logger.Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	// Start with the base logger or existing context logger
	contextLogger := z.getLoggerWithContext(ctx)

	// Create a new logger with the additional fields
	loggerWithFields := contextLogger.With()

	// Add all the provided fields using the helper method
	for _, field := range fields {
		loggerWithFields = z.addFieldToContext(loggerWithFields, field.Key(), field.Value())
	}

	// Create the final logger and embed it in the context
	finalLogger := loggerWithFields.Logger()
	return finalLogger.WithContext(ctx)
}
