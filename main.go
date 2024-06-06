package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var (
	lastCheckTime       = time.Now()
	monitoredEventNames = []string{
		"Foo Fighters",
		"Liam Gallagher",
	}
)

func main() {
	slog.Info(
		fmt.Sprintf("Monitoring: %s", strings.Join(monitoredEventNames, ", ")),
	)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	exitChan := make(chan struct{})

	for {
		select {
		case <-ticker.C:

			checkTime := time.Now()

			tickets, err := twickets.GetLatestTickets(
				context.Background(),
				twickets.GetTicketsInput{
					Country:    twickets.CountryUnitedKingdom,
					Regions:    []twickets.Region{twickets.RegionLondon},
					MaxNumber:  100,
					BeforeTime: checkTime,
				},
			)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			processTickets(tickets)

			lastCheckTime = checkTime

		case <-exitChan:
			return
		}
	}
}

func processTickets(tickets []twickets.Ticket) {
	for _, ticket := range tickets {
		if ticket.CreatedAt.Before(lastCheckTime) {
			continue
		}

		slog.Info(
			"found tickets for event",
			"event", ticket.Event.Name,
		)

		for _, eventName := range monitoredEventNames {
			isMonitored := fuzzy.MatchNormalizedFold(ticket.Event.Name, eventName)
			if isMonitored {
				slog.Info(
					"found tickets for monitored event",
					"name", ticket.Event.Name,
					"tickets", ticket.TicketQuantity,
					"ticketCost", ticket.TotalSellingPrice.PerString(ticket.TicketQuantity),
					"totalCost", ticket.TotalSellingPrice.String(),
				)
			}
		}
	}
}
