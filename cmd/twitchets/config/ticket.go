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

func (t TicketConfig) Validate() error {
	if t.Event == "" {
		return errors.New("event name must be set")
	}

	for _, region := range t.Regions {
		if !twigots.Regions.Contains(region) {
			return fmt.Errorf("region '%s' is not valid", region)
		}
	}

	for _, notification := range t.Notification {
		if !NotificationTypes.Contains(notification) {
			return fmt.Errorf("notification type '%s' is not valid", notification)
		}
	}

	return nil
}

func (t TicketConfig) Filter() twigots.Filter {
	var filter twigots.Filter
	filter.Event = t.Event

	if t.EventSimilarity != nil {
		filter.EventSimilarity = *t.EventSimilarity
	}

	if t.Regions != nil {
		filter.Regions = t.Regions
	}

	if t.NumTickets != nil {
		filter.NumTickets = *t.NumTickets
	}

	if t.Discount != nil {
		filter.MinDiscount = *t.Discount
	}

	return filter
}
