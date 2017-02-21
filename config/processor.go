package config

import (
	"bytes"

	"gopkg.in/yaml.v2"
)

// Processor handles the reading/writing of data files of different formats
type Processor interface {
	Marshal(map[string]interface{}) (*bytes.Buffer, error)
}

// YamlProcessor processes data into and out of YAML
type yamlProcessor struct{}

// NewYamlProcessor returns a processor of YAML files
func NewYamlProcessor() Processor {
	return &yamlProcessor{}
}

// Marshal data string into a yaml file
func (y *yamlProcessor) Marshal(data map[string]interface{}) (*bytes.Buffer, error) {
	m, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
