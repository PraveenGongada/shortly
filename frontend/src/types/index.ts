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

export interface SuccessResponse {
  message: string;
  success: boolean;
}

export interface CreateUrlRequest {
  long_url: string;
}

export interface CreateUrlResponse {
  id: string;
  short_url: string;
}

export interface UrlResponse {
  id: string;
  short_code: string;
  long_url: string;
  redirects: number;
}

export interface UrlUpdateRequest {
  id: string;
  new_url: string;
}

export interface DeleteUrlRequest {
  id: string;
}

export interface Token {
  type: string;
  token: string;
}

export interface UserLoginRequest {
  email: string;
  password: string;
}

export interface UserRegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface UserResponse {
  id: string;
  name: string;
  email: string;
}

export interface UserTokenResponse {
  type: string;
  token: string;
}

export interface PaginationParams {
  limit: number;
  offset: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}
