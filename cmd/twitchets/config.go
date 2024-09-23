package main

import (
	"errors"
	"fmt"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	APIKey       string            `json:"apiKey"`
	GlobalConfig GlobalEventConfig `json:"global"`
	EventConfig  []twickets.Filter `json:"events"`
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("api key must be set")
	}

	err := c.GlobalConfig.Validate()
	if err != nil {
		return fmt.Errorf("global config is not valid: %w", err)
	}

	for idx, event := range c.EventConfig {
		err := event.Validate()
		if err != nil {
			return fmt.Errorf("event config at index [%d] is no valid: %w", idx, err)
		}
	}

	return nil
}

func (c Config) applyGlobalConfig() {
	for idx, eventConfig := range c.EventConfig {
		// Apply regions
		if len(eventConfig.Regions) == 0 && len(c.GlobalConfig.Regions) != 0 {
			eventConfig.Regions = c.GlobalConfig.Regions
		}

		// Apply num tickets
		if eventConfig.NumTickets == 0 && c.GlobalConfig.NumTickets != 0 {
			eventConfig.NumTickets = c.GlobalConfig.NumTickets
		}

		// Apply discount
		if eventConfig.Discount == 0 && c.GlobalConfig.Discount != 0 {
			eventConfig.Discount = c.GlobalConfig.Discount
		}

		c.EventConfig[idx] = eventConfig
	}
}

// GlobalEventConfig is config that applies to all events,
// unless an event explicitly overwrites its.
// Country is required.
type GlobalEventConfig struct {
	Country    twickets.Country  `json:"country"`
	Regions    []twickets.Region `json:"regions"`
	NumTickets int               `json:"num_tickets"`
	Discount   float64           `json:"discount"`
}

func (c GlobalEventConfig) Validate() error {
	if c.Country.Value == "" {
		return errors.New("country must be set")
	}
	if !twickets.Countries.Contains(c.Country) {
		return fmt.Errorf("country '%s' is not valid", c.Country)
	}

	// Reuse the filter validation logic
	globalFilter := twickets.Filter{
		Name:       "global", // Name must be be set - this is arbitrary
		Regions:    c.Regions,
		NumTickets: c.NumTickets,
		Discount:   c.Discount,
	}
	err := globalFilter.Validate()
	if err != nil {
		return err
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
		"", nil,
		koanf.UnmarshalConf{
			Tag: "json",
			DecoderConfig: &mapstructure.DecoderConfig{
				// Mostly a copy of decoder config form koanf
				Result:           &config,
				WeaklyTypedInput: true,
				ErrorUnused:      true,
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					mapstructure.TextUnmarshallerHookFunc()),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	config.applyGlobalConfig()

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
