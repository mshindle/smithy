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
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type YamlProcessor struct{}

// NewYamlProcessor returns an instance of a YamlProcessor
func NewYamlProcessor() Processor {
	return &YamlProcessor{}
}

// Marshal data string into a yaml file
func (y *YamlProcessor) Marshal(data Object) (*bytes.Buffer, error) {
	m, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}

func (y *YamlProcessor) UnmarshalFile(file string) (Object, error) {
	var object Object

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return object, err
	}

	err = yaml.Unmarshal(fileBytes, &object)
	if err != nil {
		return object, err
	}

	return object, nil
}
