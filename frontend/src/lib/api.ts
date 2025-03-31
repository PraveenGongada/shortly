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

import {
  CreateUrlRequest,
  CreateUrlResponse,
  PaginatedResponse,
  PaginationParams,
  SuccessResponse,
  Token,
  UrlResponse,
  UrlUpdateRequest,
  UserLoginRequest,
  UserRegisterRequest,
  UserResponse,
} from "@/types";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.message || response.statusText || "An error occurred",
    );
  }

  return response.json() as Promise<T>;
}

// Helper to get common headers
function getCommonHeaders(): HeadersInit {
  return {
    "Content-Type": "application/json",
  };
}

// Auth related API calls
export const authApi = {
  login: async (data: UserLoginRequest): Promise<Token> => {
    const response = await fetch(`${API_BASE_URL}/user/login`, {
      method: "POST",
      headers: getCommonHeaders(),
      body: JSON.stringify(data),
      credentials: "include", // This is important for cookies
    });

    const result = await handleResponse<Token>(response);

    // Set a cookie to identify authenticated state
    document.cookie = `token=${result.token}; path=/;`;

    return result;
  },

  register: async (data: UserRegisterRequest): Promise<UserResponse> => {
    const response = await fetch(`${API_BASE_URL}/user/register`, {
      method: "POST",
      headers: getCommonHeaders(),
      body: JSON.stringify(data),
      credentials: "include",
    });

    return handleResponse<UserResponse>(response);
  },

  logout: async (): Promise<SuccessResponse> => {
    const response = await fetch(`${API_BASE_URL}/user/logout`, {
      method: "GET",
      headers: getCommonHeaders(),
      credentials: "include",
    });

    // Clear auth token cookie
    document.cookie = "token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";

    return handleResponse<SuccessResponse>(response);
  },
};

// URL shortening related API calls
export const urlApi = {
  createShortUrl: async (
    data: CreateUrlRequest,
  ): Promise<CreateUrlResponse> => {
    const response = await fetch(`${API_BASE_URL}/url/create`, {
      method: "POST",
      headers: getCommonHeaders(),
      body: JSON.stringify(data),
      credentials: "include",
    });

    return handleResponse<CreateUrlResponse>(response);
  },

  getUrls: async (
    params: PaginationParams,
  ): Promise<PaginatedResponse<UrlResponse>> => {
    const { limit, offset } = params;
    const response = await fetch(
      `${API_BASE_URL}/urls?limit=${limit}&offset=${offset}`,
      {
        method: "GET",
        headers: getCommonHeaders(),
        credentials: "include",
      },
    );

    return handleResponse<PaginatedResponse<UrlResponse>>(response);
  },

  updateUrl: async (data: UrlUpdateRequest): Promise<SuccessResponse> => {
    const response = await fetch(`${API_BASE_URL}/url/update`, {
      method: "PATCH",
      headers: getCommonHeaders(),
      body: JSON.stringify(data),
      credentials: "include",
    });

    return handleResponse<SuccessResponse>(response);
  },

  deleteUrl: async (id: string): Promise<SuccessResponse> => {
    const response = await fetch(`${API_BASE_URL}/url/${id}`, {
      method: "DELETE",
      headers: getCommonHeaders(),
      credentials: "include",
    });

    return handleResponse<SuccessResponse>(response);
  },

  getAnalytics: async (shortUrl: string): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/url/analytics/${shortUrl}`, {
      method: "GET",
      headers: getCommonHeaders(),
      credentials: "include",
    });

    return handleResponse<any>(response);
  },

  redirect: async (shortUrl: string): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/${shortUrl}`, {
      method: "GET",
      headers: getCommonHeaders(),
    });

    return handleResponse<any>(response);
  },
};
