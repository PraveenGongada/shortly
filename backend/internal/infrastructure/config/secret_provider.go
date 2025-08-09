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

package config

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type SecretProvider interface {
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetRedisPassword() string
	GetRSAPublicKey() *rsa.PublicKey
	GetRSAPrivateKey() *rsa.PrivateKey
}

type EnvSecretProvider struct {
	rsaPublicKey  *rsa.PublicKey
	rsaPrivateKey *rsa.PrivateKey
}

func NewEnvSecretProvider() (*EnvSecretProvider, error) {
	provider := &EnvSecretProvider{}

	// Validate required database credentials
	if os.Getenv("POSTGRES_USER") == "" {
		return nil, fmt.Errorf("POSTGRES_USER environment variable is required")
	}
	if os.Getenv("POSTGRES_PASSWORD") == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable is required")
	}

	// Load RSA keys
	if err := provider.loadRSAKeys(); err != nil {
		return nil, fmt.Errorf("failed to load RSA keys: %w", err)
	}

	return provider, nil
}

func (e *EnvSecretProvider) GetDatabaseUser() string {
	return os.Getenv("POSTGRES_USER")
}

func (e *EnvSecretProvider) GetDatabasePassword() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func (e *EnvSecretProvider) GetRSAPublicKey() *rsa.PublicKey {
	if e.rsaPublicKey == nil {
		panic("RSA public key not loaded - check JWT_PUBLIC_KEY_PATH configuration")
	}
	return e.rsaPublicKey
}

func (e *EnvSecretProvider) GetRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func (e *EnvSecretProvider) GetRSAPrivateKey() *rsa.PrivateKey {
	if e.rsaPrivateKey == nil {
		panic("RSA private key not loaded - check JWT_PRIVATE_KEY_PATH configuration")
	}
	return e.rsaPrivateKey
}

func (e *EnvSecretProvider) loadRSAKeys() error {
	publicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")

	if publicKeyPath == "" {
		return fmt.Errorf("JWT_PUBLIC_KEY_PATH environment variable is required")
	}
	if privateKeyPath == "" {
		return fmt.Errorf("JWT_PRIVATE_KEY_PATH environment variable is required")
	}

	publicKey, err := loadRSAPublicKeyFromFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load RSA public key from %s: %w", publicKeyPath, err)
	}
	e.rsaPublicKey = publicKey

	privateKey, err := loadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load RSA private key from %s: %w", privateKeyPath, err)
	}
	e.rsaPrivateKey = privateKey

	return nil
}

func loadRSAPublicKeyFromFile(filepath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	return publicKey, nil
}

func loadRSAPrivateKeyFromFile(filepath string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return privateKey, nil
}
