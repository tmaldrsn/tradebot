package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type SourceConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Tickers []struct {
		Ticker          string `yaml:"ticker"`
		Timeframe       string `yaml:"timeframe"`
		PollingInterval string `yaml:"polling_interval"`
	} `yaml:"tickers"`
}

type Config struct {
	Sources []SourceConfig `yaml:"sources"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	err = ValidateConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	seen := make(map[string]bool)
	for _, src := range cfg.Sources {
		if seen[src.Name] {
			return fmt.Errorf("duplicate source name found: %s", src.Name)
		}
		seen[src.Name] = true
	}
	return nil
}

func (cfg *Config) GetSource(name, typ string) (*SourceConfig, error) {
	for _, src := range cfg.Sources {
		if src.Name == name && src.Type == typ {
			return &src, nil
		}
	}
	return nil, fmt.Errorf("source config not found for name=%s, type=%s", name, typ)
}
