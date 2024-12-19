package config

import (
	"errors"

	"github.com/ahobsonsayers/twigots"
)

// GlobalEventConfig is config that applies to all events,
// unless an event explicitly overwrites its.
// Country is required.
type GlobalEventConfig struct {
	EventSimilarity float64            `json:"eventSimilarity"`
	Regions         []twigots.Region   `json:"regions"`
	NumTickets      int                `json:"numTickets"`
	Discount        float64            `json:"discount"`
	Notification    []NotificationType `json:"notification"`
}

func (c GlobalEventConfig) Validate() error {
	// Convert config to a filter
	filter := twigots.Filter{
		Event:           "global", // Event must be be set - this is arbitrary
		EventSimilarity: c.EventSimilarity,
		Regions:         c.Regions,
		NumTickets:      c.NumTickets,
		MinDiscount:     c.Discount / 100,
	}

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

	return nil
}
