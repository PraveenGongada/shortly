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

package errors

import "fmt"

// ErrorType represents different categories of domain errors
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeConflict       ErrorType = "conflict"
	ErrorTypeUnauthorized   ErrorType = "unauthorized"
	ErrorTypeForbidden      ErrorType = "forbidden"
	ErrorTypeInternal       ErrorType = "internal"
)

// DomainError represents an error in the domain layer
type DomainError struct {
	Type    ErrorType
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Domain error constructors
func ValidationError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: msg,
	}
}

func NotFoundError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: msg,
	}
}

func ConflictError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeConflict,
		Message: msg,
	}
}

func UnauthorizedError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeUnauthorized,
		Message: msg,
	}
}

func ForbiddenError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeForbidden,
		Message: msg,
	}
}

func InternalError(msg string) error {
	return &DomainError{
		Type:    ErrorTypeInternal,
		Message: msg,
	}
}

// GetErrorType extracts the error type from an error, returns ErrorTypeInternal if not a DomainError
func GetErrorType(err error) ErrorType {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type
	}
	return ErrorTypeInternal
}