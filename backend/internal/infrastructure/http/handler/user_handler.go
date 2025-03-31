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

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
	"github.com/PraveenGongada/shortly/internal/shared/utils"
)

// UserLogin godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body valueobject.UserLoginReqest true "User login credentials"
// @Success 200 {object} response.Response{data=valueobject.UserTokenRespBody} "Login successful"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/login [post]
func (h HttpHandlerImpl) UserLogin(w http.ResponseWriter, r *http.Request) {
	loginReq := &valueobject.UserLoginReqest{}
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		response.Err(w, errors.BadRequest("Error decoding Request"))
		return
	}

	res, err := h.userService.UserLogin(r.Context(), loginReq)
	if err != nil {
		response.Err(w, err)
		return
	}

	err = utils.SetCookie(w, res.Token)
	if err != nil {
		log.Err(err).Msg("Error setting the cookie")
	}

	response.Json(w, http.StatusOK, "Login successful!", res)
}

// UserLogout godoc
// @Summary User logout
// @Description Logout the current user by invalidating their cookie
// @Tags user
// @Produce json
// @Success 200 {object} response.Response "Logout successful"
// @Router /user/logout [get]
func (h HttpHandlerImpl) UserLogout(w http.ResponseWriter, r *http.Request) {
	utils.InvalidateCookie(w)
	response.Json(w, http.StatusOK, "success!", nil)
}

// UserRegsiter godoc
// @Summary User registration
// @Description Register a new user and return a JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body valueobject.UserRegisterRequest true "User registration data"
// @Success 200 {object} response.Response{data=valueobject.UserTokenRespBody} "Registration successful"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Email already registered"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/register [post]
func (h HttpHandlerImpl) UserRegsiter(w http.ResponseWriter, r *http.Request) {
	registerReq := &valueobject.UserRegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		response.Err(w, errors.BadRequest("Error decoding data"))
		return
	}

	userResponse, err := h.userService.UserRegister(r.Context(), registerReq)
	if err != nil {
		response.Err(w, err)
		return
	}

	err = utils.SetCookie(w, userResponse.Token)
	if err != nil {
		log.Err(err).Msg("Error setting the cookie")
	}

	response.Json(w, http.StatusOK, "Registration successful!", userResponse)
}
