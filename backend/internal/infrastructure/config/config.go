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
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

var (
	globalConfig  *Config
	globalSecrets SecretProvider
	configOnce    sync.Once
	secretsOnce   sync.Once
)

type Manager struct {
	loader  *Loader
	secrets SecretProvider
}

func NewManager(configPaths []string, envPrefix string) *Manager {
	return &Manager{
		loader: NewLoader(configPaths, envPrefix),
	}
}

func (m *Manager) LoadConfig() (*Config, error) {
	return m.loader.Load()
}

func (m *Manager) LoadSecrets() (SecretProvider, error) {
	if m.secrets == nil {
		var err error
		m.secrets, err = NewEnvSecretProvider()
		if err != nil {
			return nil, fmt.Errorf("failed to create secret provider: %w", err)
		}
	}
	return m.secrets, nil
}

func LoadGlobalConfig() (*Config, error) {
	var err error
	configOnce.Do(func() {
		_ = godotenv.Load()

		manager := NewManager([]string{".", "/etc/shortly/", "configs/"}, "")
		globalConfig, err = manager.LoadConfig()
	})
	return globalConfig, err
}

func LoadGlobalSecrets() (SecretProvider, error) {
	var err error
	secretsOnce.Do(func() {
		_ = godotenv.Load()

		globalSecrets, err = NewEnvSecretProvider()
	})
	return globalSecrets, err
}

func GetGlobalConfig() *Config {
	config, err := LoadGlobalConfig()
	if err != nil {
		panic("Failed to load global configuration: " + err.Error())
	}
	return config
}

func GetGlobalSecrets() SecretProvider {
	secrets, err := LoadGlobalSecrets()
	if err != nil {
		panic("Failed to load global secrets: " + err.Error())
	}
	return secrets
}
