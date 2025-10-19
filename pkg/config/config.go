package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func GetConfig[T any](reader io.Reader) (*T, error) {
	var cfg T
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("the config content is malformed: %w", err)
	}

	return &cfg, nil
}

func LoadYAMLDocument[T any](path string) (*T, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return GetConfig[T](bytes.NewReader(b))
}
