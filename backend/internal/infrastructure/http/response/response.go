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

package response

import (
	"encoding/json"
	"net/http"

	"github.com/PraveenGongada/shortly/internal/shared/errors"
)

type Response struct {
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Json(w http.ResponseWriter, httpCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	res := Response{
		Message: &message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(res)
}

func Text(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpCode)
	w.Write([]byte(message))
}

func Err(w http.ResponseWriter, err error) {
	_, ok := err.(*errors.ErrorResp)
	if !ok {
		err = errors.InternalServerError()
	}

	er, _ := err.(*errors.ErrorResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(er.Code)
	res := Response{
		Message: &er.Message,
	}
	json.NewEncoder(w).Encode(res)
}
