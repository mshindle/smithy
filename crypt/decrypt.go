// Copyright © 2017 Michael Shindle <mshindle@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
)

// DecryptFromString decrypts a standard base64
// encoded string, as defined in RFC 4648, wrapped in an
// "ENC[" and "]" construct
func DecryptFromString(s string, label string, file string) ([]byte, error) {
	b64string := s[4 : len(s)-1]
	log.WithField("base64", b64string).Debug("string to decode")
	decodeBytes, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		return nil, err
	}
	return Decrypt(decodeBytes, label, file)
}

// Decrypt will decrypt the data bytes using the PEM
// private key
func Decrypt(data []byte, label string, file string) ([]byte, error) {
	var privkey rsa.PrivateKey
	err := loadKey(file, &privkey)
	if err != nil {
		log.WithField("file", file).WithError(err).Error("could not load private key")
		return nil, err
	}

	decryptedValue, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &privkey, data, []byte(label))
	if err != nil {
		log.WithError(err).WithField("value", decryptedValue).Error("could not decrypt data")
		return nil, err
	}

	return decryptedValue, nil
}
