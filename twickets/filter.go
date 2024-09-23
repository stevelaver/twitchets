package twickets // nolint

import (
	"errors"
	"fmt"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/samber/lo"
)

type Filter struct {
	Name       string   `json:"name"`
	Regions    []Region `json:"regions"`
	NumTickets int      `json:"num_tickets"`
	Discount   float64  `json:"discount"`
}

func (f Filter) Validate() error {
	if f.Name == "" {
		return errors.New("event name must be set")
	}

	for _, region := range f.Regions {
		if !Regions.Contains(region) {
			return fmt.Errorf("region '%s' is not valid", region)
		}
	}

	if f.NumTickets < 0 {
		return errors.New("number of tickets cannot be negative")
	}

	if f.Discount < 0 {
		return errors.New("discount cannot be negative")
	}

	return nil
}

// TicketMatches check is a ticket matches the filter
func (f Filter) TicketMatches(ticket Ticket) bool {
	return matchesEventName(ticket, f.Name) &&
		matchesRegions(ticket, f.Regions) &&
		matchesNumTickets(ticket, f.NumTickets) &&
		matchesDiscount(ticket, f.Discount)
}

// matchesEventName returns whether a tickets matches a desired event name
func matchesEventName(ticket Ticket, eventName string) bool {
	ticketEventName := normaliseEventName(ticket.Event.Name)
	desiredEventName := normaliseEventName(eventName)

	similarity := strutil.Similarity(ticketEventName, desiredEventName, metrics.NewJaroWinkler())
	return similarity >= 0.85
}

// matchesRegions determines whether a tickets matches any of desired regions
func matchesRegions(ticket Ticket, regions []Region) bool {
	if len(regions) == 0 {
		return true
	}
	return lo.Contains(regions, ticket.Event.Venue.Location.RegionCode)
}

// matchesNumTickets determines whether a tickets matches a desired number of tickets
func matchesNumTickets(ticket Ticket, numTickets int) bool {
	if numTickets <= 0 {
		return true
	}
	return ticket.TicketQuantity == numTickets
}

// matchesDiscount determines whether a tickets matches a desired discount
func matchesDiscount(ticket Ticket, discount float64) bool {
	if discount <= 0 {
		return true
	}
	return ticket.Discount() >= discount
}
