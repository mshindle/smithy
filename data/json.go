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

	"encoding/json"
	"io/ioutil"
)

// JsonProcessor handles transforming Objects into JSON and vice-versa
type JsonProcessor struct{}

// NewJsonProcessor returns a JSON backed processor.
func NewJsonProcessor() Processor {
	return &JsonProcessor{}
}

func (j *JsonProcessor) Marshal(data Object) (*bytes.Buffer, error) {
	m, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}

// UnmarshalFile will read the contents of file, unmarshal the json, and convert any encrypted data fields
func (j *JsonProcessor) UnmarshalFile(file string) (Object, error) {
	var object Object

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return object, err
	}

	err = json.Unmarshal(fileBytes, &object)
	if err != nil {
		return object, err
	}

	return object, nil
}
