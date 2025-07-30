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

package repository

import (
	"context"

	"github.com/PraveenGongada/shortly/internal/domain/url/entity"
)

// URLRepository defines persistence operations for URLs
type URLRepository interface {
	Save(ctx context.Context, url *entity.URL) (*entity.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (*entity.URL, error)
	FindByID(ctx context.Context, id string) (*entity.URL, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.URL, error)
	ExistsByShortCode(ctx context.Context, shortCode string) (bool, error)
	Update(ctx context.Context, url *entity.URL) error
	Delete(ctx context.Context, id, userID string) error
	IncrementRedirects(ctx context.Context, shortCode string) error
}
