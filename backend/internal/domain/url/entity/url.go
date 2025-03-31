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

package entity

import "github.com/guregu/null/v5"

type Url struct {
	Id        string      `db:"id"`
	UserID    string      `db:"user_id"`
	ShortUrl  string      `db:"short_url"`
	LongUrl   string      `db:"long_url"`
	Redirects int         `db:"redirects"`
	CreatedAt string      `db:"created_at"`
	UpdatedAt null.String `db:"updated_at"`
}
