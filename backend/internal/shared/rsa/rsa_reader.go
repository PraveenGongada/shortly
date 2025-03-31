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

package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ReadPublicKeyFromEnv(rsaPublic string) (*rsa.PublicKey, error) {
	data, _ := pem.Decode([]byte(rsaPublic))
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		return nil, errors.New("cannot reflect the interface")
	}

	return publicKey, nil
}

func ReadPrivateKeyFromEnv(rsaPrivate string) (any, error) {
	data, _ := pem.Decode([]byte(rsaPrivate))
	privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKeyImported, nil
}
