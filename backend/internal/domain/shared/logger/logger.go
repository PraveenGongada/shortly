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

package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	WithContext(ctx context.Context, fields ...Field) context.Context
}

// Field represents a structured logging field
type Field interface {
	Key() string
	Value() interface{}
}

type field struct {
	key   string
	value interface{}
}

func (f field) Key() string {
	return f.key
}

func (f field) Value() interface{} {
	return f.value
}

func String(key, value string) Field {
	return field{key: key, value: value}
}

func Int(key string, value int) Field {
	return field{key: key, value: value}
}

func Int64(key string, value int64) Field {
	return field{key: key, value: value}
}

func Float64(key string, value float64) Field {
	return field{key: key, value: value}
}

func Bool(key string, value bool) Field {
	return field{key: key, value: value}
}

func Error(err error) Field {
	return field{key: "error", value: err}
}

func Any(key string, value interface{}) Field {
	return field{key: key, value: value}
}