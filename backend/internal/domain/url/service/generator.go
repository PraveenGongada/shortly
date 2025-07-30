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

package service

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/PraveenGongada/shortly/internal/domain/interfaces"
)

const (
	DefaultShortCodeLength = 6
	shortCodeCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type generator struct {
	shortCodeLength int
}

func NewGenerator(shortCodeLength int) interfaces.ShortCodeGenerator {
	if shortCodeLength <= 0 {
		shortCodeLength = DefaultShortCodeLength
	}
	return &generator{
		shortCodeLength: shortCodeLength,
	}
}

func (g *generator) GenerateShortCode() (string, error) {
	result := make([]byte, g.shortCodeLength)
	charsetLength := big.NewInt(int64(len(shortCodeCharset)))

	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", errors.New("failed to generate random short code")
		}
		result[i] = shortCodeCharset[randomIndex.Int64()]
	}

	return string(result), nil
}
