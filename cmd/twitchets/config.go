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
	Regions []twickets.Region `json:"regions"`
	Events  []twickets.Filter `json:"events"`
}

func (c *Config) parseKoanf(k *koanf.Koanf) error {
	if k == nil {
		return nil
	}

	// Parse country
	countryString := k.String("country")
	if countryString == "" {
		return errors.New("country must be set")
	}
	country := twickets.Countries.Parse(countryString)
	if country == nil {
		return fmt.Errorf("%s is not a valid country code", countryString)
	}

	// Parse regions
	regionStrings := k.Strings("regions")
	regions := make([]twickets.Region, 0, len(regionStrings))
	for _, regionString := range regionStrings {
		region := twickets.Regions.Parse(regionString)
		if region == nil {
			return fmt.Errorf("%s is not a valid region code", countryString)
		}
		regions = append(regions, *region)
	}

	// Parse filters
	var events []twickets.Filter
	err := k.UnmarshalWithConf(
		"events", &events,
		koanf.UnmarshalConf{Tag: `json`},
	)
	if err != nil {
		return fmt.Errorf("invalid events: %w", err)
	}

	c.Country = *country
	c.Regions = regions
	c.Events = events

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
