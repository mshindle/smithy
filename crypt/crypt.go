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

	log "github.com/Sirupsen/logrus"
)

const BitSize = 1024

// Generate creates a public / private key pair and saves them in the specified file
func Generate(pubfile string, privfile string) error {
	key, err := rsa.GenerateKey(rand.Reader, BitSize)
	if err != nil {
		return err
	}
	pubkey := key.PublicKey

	log.WithFields(log.Fields{
		"prime0":         key.Primes[0].String(),
		"prime1":         key.Primes[1].String(),
		"exponent":       key.D.String(),
		"pubkeyModulus":  pubkey.N.String(),
		"pubkeyExponent": pubkey.E,
	}).Info("key generated")

	err = saveKey(privfile, key)
	if err != nil {
		log.WithField("file", privfile).Error("could not save private key")
		return err
	}

	err = saveKey(pubfile, pubkey)
	if err != nil {
		log.WithField("file", pubfile).Error("could not save public key")
		return err
	}

	log.WithFields(log.Fields{
		"privfile": privfile,
		"pubfile":  pubfile,
	}).Info("all key files written")

	return nil
}
