package config

import (
	"errors"
	"fmt"

	"github.com/ahobsonsayers/twitchets/twickets"
)

type Config struct {
	APIKey        string            `json:"apiKey"`
	Country       twickets.Country  `json:"country"`
	GlobalConfig  GlobalEventConfig `json:"global"`
	TicketsConfig []TicketConfig    `json:"tickets"`
}

func (c Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("api key must be set")
	}

	if c.Country.Value == "" {
		return errors.New("country must be set")
	}
	if !twickets.Countries.Contains(c.Country) {
		return fmt.Errorf("country '%s' is not valid", c.Country)
	}

	err := c.GlobalConfig.Validate()
	if err != nil {
		return fmt.Errorf("global config is not valid: %w", err)
	}

	for idx, ticketConfig := range c.TicketsConfig {
		err := ticketConfig.Validate()
		if err != nil {
			return fmt.Errorf("event config at index [%d] is no valid: %w", idx, err)
		}
	}

	return nil
}

func (c Config) Filters() []twickets.Filter {
	filters := make([]twickets.Filter, 0, len(c.TicketsConfig))
	for _, ticketConfig := range c.TicketsConfig {

		var filter twickets.Filter
		filter.Event = ticketConfig.Event

		// Set name similarity
		if ticketConfig.EventSimilarity == nil {
			filter.EventSimilarity = c.GlobalConfig.EventSimilarity
		} else if *ticketConfig.EventSimilarity > 0 {
			filter.EventSimilarity = *ticketConfig.EventSimilarity
		}

		// Set regions
		if ticketConfig.Regions == nil {
			filter.Regions = c.GlobalConfig.Regions
		} else {
			filter.Regions = ticketConfig.Regions
		}

		// Set num tickets
		if ticketConfig.NumTickets == nil {
			filter.NumTickets = c.GlobalConfig.NumTickets
		} else if *ticketConfig.NumTickets > 0 {
			filter.NumTickets = *ticketConfig.NumTickets
		}

		// Set discount
		if ticketConfig.Discount == nil {
			filter.Discount = c.GlobalConfig.Discount
		} else if *ticketConfig.Discount > 0 {
			filter.Discount = *ticketConfig.Discount
		}

		filters = append(filters, filter)
	}

	return filters
}
