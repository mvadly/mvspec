package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title       string         `yaml:"title"`
	Version     string         `yaml:"version"`
	Description string         `yaml:"description"`
	Host        string         `yaml:"host"`
	BasePath    string         `yaml:"basePath"`
	Output      string         `yaml:"output"`
	Exclude     []string       `yaml:"exclude"`
	ParseTypes  bool           `yaml:"parseTypes"`
	Servers     []ServerConfig `yaml:"servers"`
	EnvVars     []EnvVarConfig `yaml:"env"`
}

type ServerConfig struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

type EnvVarConfig struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Output == "" {
		cfg.Output = "mv-spec.json"
	}

	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
