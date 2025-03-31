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

import (
	"fmt"
	"net/http"
)

type ErrorResp struct {
	Code    int
	Message string
}

func (r *ErrorResp) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func InternalServerError() error {
	return &ErrorResp{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong!",
	}
}

func BadRequest(msg string) error {
	return &ErrorResp{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NotFound(msg string) error {
	return &ErrorResp{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func Unauthorized(msg string) error {
	return &ErrorResp{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}
