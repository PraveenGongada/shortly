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
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/PraveenGongada/shortly/internal/domain/url/valueobject"
	"github.com/PraveenGongada/shortly/internal/infrastructure/http/response"
	"github.com/PraveenGongada/shortly/internal/shared/errors"
)

// CreateShortUrl godoc
// @Summary Create a short URL
// @Description Create a short URL from a long URL
// @Tags url
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body valueobject.CreateUrlRequest true "URL information"
// @Success 201 {object} response.Response{data=valueobject.CreateUrlResponse} "URL created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/create [post]
func (h HttpHandlerImpl) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context()).With().Str("handler", "CreateShortUrl").Logger()

	req := &valueobject.CreateUrlRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	validUrl := strings.TrimSpace(req.LongUrl) != "" || strings.HasPrefix(req.LongUrl, "http") ||
		strings.HasPrefix(req.LongUrl, "https")

	if err != nil || !validUrl {
		logger.Warn().Err(err).Str("longUrl", req.LongUrl).Msg("Invalid request payload")
		response.Err(w, errors.BadRequest("Invalid Request"))
		return
	}

	userId := r.Header.Get("id")

	urlResponse, err := h.urlService.CreateShortUrl(r.Context(), userId, req)
	if err != nil {
		logger.Error().Err(err).Msg("Service layer error creating short URL")
		response.Err(w, err)
		return
	}

	logger.Info().
		Str("shortUrl", urlResponse.ShortUrl).
		Str("urlId", urlResponse.Id).
		Msg("Short URL creation request completed successfully")

	response.Json(w, http.StatusCreated, "Short url created successfully", urlResponse)
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
func (h HttpHandlerImpl) RedirectUser(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortUrl")

	logger := log.Ctx(r.Context()).
		With().
		Str("handler", "RedirectUser").
		Str("shortUrl", shortUrl).
		Logger()

	if shortUrl == "favicon.ico" {
		return
	}

	longUrl, err := h.urlService.GetLongUrl(r.Context(), shortUrl)
	if err != nil {
		logger.Warn().Err(err).Msg("Redirect failed - URL not found or error")
		response.Err(w, err)
		return
	}

	logger.Info().Str("longUrl", longUrl).Msg("Redirect successful")
	http.Redirect(w, r, longUrl, http.StatusFound)
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
func (h HttpHandlerImpl) GetLongUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortUrl")

	longUrl, err := h.urlService.GetLongUrl(r.Context(), shortUrl)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Json(w, http.StatusOK, "success!", longUrl)
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
func (h HttpHandlerImpl) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortUrl")
	userId := r.Header.Get("id")

	count, err := h.urlService.GetAnalytics(r.Context(), shortUrl, userId)
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
// @Success 200 {object} response.Response{data=[]valueobject.UrlResponse} "URLs list"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /urls [get]
func (h HttpHandlerImpl) GetPaginatedUrls(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if offsetStr == "" || limitStr == "" {
		response.Err(w, errors.BadRequest("limit & offset are required"))
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.Err(w, errors.BadRequest("error parsing offset"))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.Err(w, errors.BadRequest("error parsing limit"))
		return
	}

	urls, err := h.urlService.GetPaginatedUrls(r.Context(), userId, limit, offset)
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
// @Param request body valueobject.UrlUpdateRequest true "URL update information"
// @Success 200 {object} response.Response "URL updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /url/update [patch]
func (h HttpHandlerImpl) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	var updateRequest valueobject.UrlUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		log.Err(err).Msg("Error")
		response.Err(w, errors.BadRequest("Cannot parse update request"))
		return
	}

	err = h.urlService.UpdateUrl(r.Context(), updateRequest.Id, updateRequest.Url)
	if err != nil {
		response.Err(w, err)
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
func (h HttpHandlerImpl) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "urlId")
	userId := r.Header.Get("id")

	err := h.urlService.DeleteUrl(r.Context(), urlId, userId)
	if err != nil {
		response.Err(w, err)
	}

	response.Json(w, http.StatusOK, "success!", nil)
}
