package twickets // nolint

import (
	"errors"

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

func (t Filter) validate() error {
	if t.Name == "" {
		return errors.New("event name must be set")
	}
	return nil
}

// TicketMatches check is a ticket matches the filter
func (t Filter) TicketMatches(ticket Ticket) bool {
	return matchesEventName(ticket, t.Name) &&
		matchesRegions(ticket, t.Regions) &&
		matchesNumTickets(ticket, t.NumTickets) &&
		matchesDiscount(ticket, t.Discount)
}

// matchesEventName returns whether a tickets matches a desired event name
func matchesEventName(ticket Ticket, eventName string) bool {
	ticketEventName := normaliseEventName(ticket.Event.Name)
	desiredEventName := normaliseEventName(eventName)

	similarity := strutil.Similarity(ticketEventName, desiredEventName, metrics.NewJaroWinkler())
	if similarity >= 0.85 {
		return true
	}

	return false
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
