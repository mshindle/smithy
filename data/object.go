// Copyright © 2017 Michael Shindle <mshindle@riotgames.com>
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

package data

import (
	"bytes"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/mshindle/smithy/crypt"
)

// Object represent a structured data map with string keys
type Object map[string]interface{}

// Processor manages transforming objects into specific data formats
type Processor interface {
	Marshal(Object) (*bytes.Buffer, error)
	UnmarshalFile(string) (Object, error)
}

// DecryptValue decrypts encrypted values in an object
func (object Object) DecryptValues(label string, file string) error {
	match, err := regexp.Compile("^ENC\\[*")
	if err != nil {
		log.WithError(err).Error("cannot compile ENC regex pattern")
		return err
	}

	return object.decrypt(match, label, file)
}

func (object Object) decrypt(match *regexp.Regexp, label string, file string) error {
	for k, v := range object {
		s, ok := v.(string)
		if ok && match.MatchString(s) {
			b, err := crypt.DecryptFromString(v.(string), label, file)
			if err != nil {
				return err
			}
			object[k] = string(b)
		} else if m, ok := v.(map[string]interface{}); ok {
			log.WithField("key", k).Debug("looping on key")
			err := Object(m).decrypt(match, label, file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
