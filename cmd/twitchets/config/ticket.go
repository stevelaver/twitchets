package config

import (
	"errors"
	"fmt"

	"github.com/ahobsonsayers/twitchets/twickets"
)

type TicketConfig struct {
	Event           string            `json:"event"`
	EventSimilarity *float64          `json:"eventSimilarity"`
	Regions         []twickets.Region `json:"regions"`
	NumTickets      *int              `json:"numTickets"`
	Discount        *float64          `json:"discount"`
}

func (t TicketConfig) Validate() error {
	if t.Event == "" {
		return errors.New("event name must be set")
	}

	for _, region := range t.Regions {
		if !twickets.Regions.Contains(region) {
			return fmt.Errorf("region '%s' is not valid", region)
		}
	}

	return nil
}
