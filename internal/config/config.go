package config

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

func Default() *Config {
	return &Config{
		Logger: &Logger{
			Level:  "info",
			Format: "text",
		},
		Server: &Server{
			Host: "localhost",
			Port: "8080",
		},
	}
}

func Parse(r io.Reader) (*Config, error) {
	cfg := Default()
	if err := yaml.NewDecoder(r).Decode(cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}

	return cfg, nil
}

func ParseFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	defer file.Close()

	res, err := Parse(file)
	if err != nil {
		return nil, err
	}

	return res, nil
}
