package config

import (
	"errors"
	"fmt"

	"github.com/ahobsonsayers/twigots"
)

type TicketConfig struct {
	Event           string             `json:"event"`
	EventSimilarity *float64           `json:"eventSimilarity"`
	Regions         []twigots.Region   `json:"regions"`
	NumTickets      *int               `json:"numTickets"`
	Discount        *float64           `json:"discount"`
	Notification    []NotificationType `json:"notification"`
}

func (c TicketConfig) Validate() error {
	// Convert config to a filter
	filter := c.Filter()

	// Return a specific error for a discount above the maximum.
	// We do this as twigots returns a error specifying not > 1, which is misleading
	if filter.MinDiscount > 1 {
		return errors.New("discount cannot be > 100")
	}

	// Validate filter
	err := filter.Validate()
	if err != nil {
		return err
	}

	for _, notification := range c.Notification {
		if !NotificationTypes.Contains(notification) {
			return fmt.Errorf("notification type '%s' is not valid", notification)
		}
	}

	return nil
}

func (c TicketConfig) Filter() twigots.Filter {
	var filter twigots.Filter
	filter.Event = c.Event

	if c.EventSimilarity != nil {
		filter.EventSimilarity = *c.EventSimilarity
	}

	if c.Regions != nil {
		filter.Regions = c.Regions
	}

	if c.NumTickets != nil {
		filter.NumTickets = *c.NumTickets
	}

	if c.Discount != nil {
		filter.MinDiscount = *c.Discount / 100
	}

	return filter
}
