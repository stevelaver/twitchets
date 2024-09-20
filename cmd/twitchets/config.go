package main

import (
	"errors"
	"fmt"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Country twickets.Country  `json:"country"`
	Events  []twickets.Filter `json:"events"`
}

func (c Config) Validate() error {
	if c.Country.Value == "" {
		return errors.New("country must be set")
	}

	for _, event := range c.Events {
		err := event.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) parseKoanf(k *koanf.Koanf) error {
	if k == nil {
		return nil
	}

	// Parse config
	var config Config
	err := k.UnmarshalWithConf(
		"", &config,
		koanf.UnmarshalConf{Tag: `json`},
	)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	*c = config
	return nil
}

func LoadConfig(filePath string) (Config, error) {
	// Load config.
	k := koanf.New(".")
	err := k.Load(file.Provider(filePath), yaml.Parser())
	if err != nil {
		return Config{}, fmt.Errorf("error loading config: %w", err)
	}

	// Parse config
	config := Config{}
	err = config.parseKoanf(k)
	if err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
