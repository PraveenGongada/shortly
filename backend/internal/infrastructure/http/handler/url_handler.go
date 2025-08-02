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
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/PraveenGongada/shortly/internal/domain/shared/logger"
	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/domain/shared/errors"
)

// CreateShortUrl godoc
// @Summary Create a short URL
// @Description Create a short URL from a long URL
// @Tags url
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body valueobject.CreateURLRequest true "URL information"
// @Success 201 {object} response.Response{data=valueobject.CreateURLResponse} "URL created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/create [post]
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req valueobject.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(r.Context(), "Invalid request payload",
			logger.String("handler", "CreateShortURL"),
			logger.Error(err))
		response.Err(w, errors.ValidationError("Invalid request format"))
		return
	}

	userID := r.Header.Get("id")

	h.logger.Info(r.Context(), "Processing create short URL request",
		logger.String("handler", "CreateShortURL"),
		logger.String("userID", userID),
		logger.String("longURL", req.LongURL))

	urlResponse, err := h.urlService.CreateShortURL(r.Context(), userID, &req)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusCreated, "Short URL created successfully", urlResponse)
}

// RedirectUser godoc
// @Summary Redirect to long URL
// @Description Redirects to the original long URL from a short URL
// @Tags url
// @Produce html
// @Param shortUrl path string true "Short URL code"
// @Success 302 {string} string "Redirect to long URL"
// @Failure 404 {object} response.Response "URL not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /{shortUrl} [get]
func (h *Handler) RedirectUser(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortUrl")

	if shortCode == "favicon.ico" {
		return
	}

	h.logger.Info(r.Context(), "Processing redirect request",
		logger.String("handler", "RedirectUser"),
		logger.String("shortCode", shortCode))

	longURL, err := h.urlService.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		response.Err(w, err)
		return
	}

	h.logger.Info(r.Context(), "Redirect successful",
		logger.String("handler", "RedirectUser"),
		logger.String("shortCode", shortCode))
	http.Redirect(w, r, longURL, http.StatusFound)
}

// GetLongUrl godoc
// @Summary Get long URL
// @Description Get the original long URL from a short URL without redirecting
// @Tags url
// @Produce json
// @Param shortUrl path string true "Short URL code"
// @Success 200 {object} response.Response{data=string} "Long URL"
// @Failure 404 {object} response.Response "URL not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /{shortUrl} [get]
func (h *Handler) GetLongURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortUrl")

	h.logger.Info(r.Context(), "Processing get long URL request",
		logger.String("handler", "GetLongURL"),
		logger.String("shortCode", shortCode))

	longURL, err := h.urlService.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusOK, "success!", longURL)
}

// GetAnalytics godoc
// @Summary Get URL analytics
// @Description Get analytics for a specific short URL
// @Tags url
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shortUrl path string true "Short URL code"
// @Success 200 {object} response.Response{data=int} "Redirect count"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "URL not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/analytics/{shortUrl} [get]
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortUrl")
	userID := r.Header.Get("id")

	h.logger.Info(r.Context(), "Processing analytics request",
		logger.String("handler", "GetAnalytics"),
		logger.String("shortCode", shortCode),
		logger.String("userID", userID))

	count, err := h.urlService.GetAnalytics(r.Context(), shortCode, userID)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusOK, "success!", count)
}

// GetPaginatedUrls godoc
// @Summary Get paginated URLs
// @Description Get a paginated list of URLs created by the user
// @Tags url
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {object} response.Response{data=[]valueobject.URLResponse} "URLs list"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /urls [get]
func (h *Handler) GetPaginatedURLs(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if offsetStr == "" || limitStr == "" {
		response.Err(w, errors.ValidationError("limit & offset are required"))
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.Err(w, errors.ValidationError("error parsing offset"))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Err(w, errors.ValidationError("error parsing limit"))
		return
	}

	urls, err := h.urlService.GetPaginatedURLs(r.Context(), userID, limit, offset)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.Json(w, http.StatusOK, "success!", urls)
}

// UpdateUrl godoc
// @Summary Update URL
// @Description Update a long URL
// @Tags url
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body valueobject.URLUpdateRequest true "URL update information"
// @Success 200 {object} response.Response "URL updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/update [patch]
func (h *Handler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	var updateRequest valueobject.URLUpdateRequest
	userID := r.Header.Get("id")

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		response.Err(w, errors.ValidationError("Cannot parse update request"))
		return
	}

	err := h.urlService.UpdateURL(r.Context(), updateRequest.ID, userID, updateRequest.NewURL)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusOK, "success!", nil)
}

// DeleteUrl godoc
// @Summary Delete URL
// @Description Delete a URL
// @Tags url
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param urlId path string true "URL ID"
// @Success 200 {object} response.Response "URL deleted successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/{urlId} [delete]
func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "urlId")
	userID := r.Header.Get("id")

	err := h.urlService.DeleteURL(r.Context(), urlID, userID)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusOK, "URL deleted successfully!", nil)
}
