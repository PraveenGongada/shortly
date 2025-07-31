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

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/user/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

// UserLogin godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body valueobject.LoginRequest true "User login credentials"
// @Success 200 {object} response.Response{data=valueobject.TokenResponse} "Login successful"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/login [post]
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq valueobject.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		response.Err(w, errors.ValidationError("Invalid request format"))
		return
	}

	tokenResp, err := h.userService.Login(r.Context(), &loginReq)
	if err != nil {
		response.Err(w, err)
		return
	}

	// Set authentication cookie
	if err := h.cookieManager.SetAuthCookie(w, tokenResp.Token); err != nil {
		h.logger.Warn(r.Context(), "Error setting login cookie", logger.Error(err))
	}

	response.Json(w, http.StatusOK, "Login successful!", tokenResp)
}

// UserLogout godoc
// @Summary User logout
// @Description Logout the current user by invalidating their cookie
// @Tags user
// @Produce json
// @Success 200 {object} response.Response "Logout successful"
// @Router /user/logout [get]
func (h *Handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	// Call service logout (for potential token blacklisting)
	if err := h.userService.Logout(r.Context()); err != nil {
		h.logger.Warn(r.Context(), "Error during logout", logger.Error(err))
	}

	// Invalidate cookie
	h.cookieManager.InvalidateAuthCookie(w)
	response.Json(w, http.StatusOK, "Logout successful!", nil)
}

// UserRegsiter godoc
// @Summary User registration
// @Description Register a new user and return a JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body valueobject.RegisterRequest true "User registration data"
// @Success 200 {object} response.Response{data=valueobject.TokenResponse} "Registration successful"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Email already registered"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/register [post]
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerReq valueobject.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		response.Err(w, errors.ValidationError("Invalid request format"))
		return
	}

	tokenResp, err := h.userService.Register(r.Context(), &registerReq)
	if err != nil {
		response.Err(w, err)
		return
	}

	// Set authentication cookie
	if err := h.cookieManager.SetAuthCookie(w, tokenResp.Token); err != nil {
		h.logger.Warn(r.Context(), "Error setting registration cookie", logger.Error(err))
	}

	response.Json(w, http.StatusCreated, "Registration successful!", tokenResp)
}
