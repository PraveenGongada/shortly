/*
 * Copyright 2026 Praveen Kumar
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

import type {
  ApiEnvelope,
  CreateUrlRequest,
  CreateUrlResponse,
  LoginRequest,
  RegisterRequest,
  TokenResponse,
  UrlResponse,
  UrlUpdateRequest,
} from "@/types";

const API_BASE = (import.meta.env.VITE_BACKEND_URL || "/api").replace(/\/+$/, "");

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

/**
 * Turn an unknown thrown value into a user-facing message. `byStatus` maps
 * specific HTTP statuses to friendly copy; otherwise the server message is
 * used for ApiErrors, and `fallback` for anything else.
 */
export function apiErrorMessage(
  err: unknown,
  fallback: string,
  byStatus?: Record<number, string>,
): string {
  if (err instanceof ApiError) {
    return byStatus?.[err.status] ?? err.message;
  }
  return fallback;
}

type RequestOptions = {
  method?: string;
  body?: unknown;
};

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetch(`${API_BASE}${path}`, {
      method: options.method ?? "GET",
      credentials: "include",
      headers:
        options.body !== undefined
          ? { "Content-Type": "application/json" }
          : undefined,
      body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
    });
  } catch {
    throw new ApiError("Network error — is the server reachable?", 0);
  }

  const raw = await res.text();
  let envelope: ApiEnvelope<T> | null = null;
  if (raw) {
    try {
      envelope = JSON.parse(raw) as ApiEnvelope<T>;
    } catch {
      envelope = null;
    }
  }

  if (!res.ok) {
    throw new ApiError(
      envelope?.message || res.statusText || "Something went wrong",
      res.status,
    );
  }

  return envelope?.data as T;
}

export const authApi = {
  login: (data: LoginRequest) =>
    request<TokenResponse>("/user/login", { method: "POST", body: data }),

  register: (data: RegisterRequest) =>
    request<TokenResponse>("/user/register", { method: "POST", body: data }),

  logout: () => request<null>("/user/logout", { method: "GET" }),
};

export const urlApi = {
  create: (data: CreateUrlRequest) =>
    request<CreateUrlResponse>("/url/create", { method: "POST", body: data }),

  list: (limit: number, offset: number) =>
    request<UrlResponse[] | null>(`/urls?limit=${limit}&offset=${offset}`),

  update: (data: UrlUpdateRequest) =>
    request<null>("/url/update", { method: "PATCH", body: data }),

  remove: (id: string) =>
    request<null>(`/url/${id}`, { method: "DELETE" }),

  resolve: (code: string) => request<string>(`/${encodeURIComponent(code)}`),
};
