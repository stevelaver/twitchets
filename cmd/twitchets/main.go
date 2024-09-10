package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/ahobsonsayers/twitchets/twickets/notification"
	"github.com/joho/godotenv"
)

const (
	maxNumTickets = 250
	refetchTime   = 1 * time.Minute
)

var (
	// Config variables
	// NOTE:
	// Region coeds are currently ignored due to
	// issues with the twickets filter api
	countryCode         = "GB"
	regionCodes         = []string{"GBLO"} // TODO reenable. See note in config variables.
	monitoredEventNames = []string{
		// Theatre
		"Back to the Future",
		"Frozen",
		"Hadestown",
		"Hamilton",
		"Harry Potter & the Cursed Child",
		"Kiss Me Kate",
		"Lion King",
		"Matilda",
		"Mean Girls",
		"Moulin Rouge",
		"My Neighbour Totoro",
		"Operation Mincemeat",
		"Starlight Express",
		"Stranger Things",
		"The Phantom Opera",
		"The Wizard of Oz",
		// Gigs
		"Coldplay",
		"Gary Clark Jr.",
		"Glass Animals",
		"Jungle",
		"Oasis",
		"Taylor Swift",
	}

	// Package variables
	country       twickets.Country
	regions       []twickets.Region
	lastCheckTime = time.Now()
)

func init() {
	_ = godotenv.Load()
}

func main() {
	// Twickets client
	parsedCountryCode := twickets.Countries.Parse(countryCode)
	if parsedCountryCode == nil {
		log.Fatalf("'%s' is not a valid country code", parsedCountryCode)
	}
	country = *parsedCountryCode

	regions = make([]twickets.Region, 0, len(regionCodes))
	for _, regionCode := range regionCodes {
		parsedRegionCode := twickets.Regions.Parse(regionCode)
		if parsedRegionCode == nil {
			log.Fatalf("'%s' is not a valid region code", parsedRegionCode)
		}
		regions = append(regions, *parsedRegionCode)
	}

	twicketsClient := twickets.NewClient(nil)

	// Notification Client
	notificationClient, err := notification.NewNtfyClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	slog.Info(
		fmt.Sprintf("Monitoring: %s", strings.Join(monitoredEventNames, ", ")),
	)

	// Initial execution
	fetchAndProcessTickets(twicketsClient, notificationClient)

	// Create ticker
	ticker := time.NewTicker(refetchTime)
	defer ticker.Stop()

	// Loop until exit
	exitChan := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			fetchAndProcessTickets(twicketsClient, notificationClient)
		case <-exitChan:
			return
		}
	}
}

func fetchAndProcessTickets(
	twicketsClient *twickets.Client,
	notificationClient notification.Client,
) {
	checkTime := time.Now()
	defer func() {
		lastCheckTime = checkTime
	}()

	tickets, err := twicketsClient.FetchTickets(
		context.Background(),
		twickets.FetchTicketsInput{
			Country: country,
			// Regions:    regions, // TODO reenable. See note in config variables.
			CreatedBefore: time.Now(),
			CreatedAfter:  lastCheckTime,
			MaxTickets:    maxNumTickets,
		},
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(tickets) == maxNumTickets {
		slog.Warn("Fetched the max number of tickets allowed. It is possible tickets have been missed.")
	}

	filteredTickets := tickets.Filter(
		twickets.TicketFilter{
			EventNames: monitoredEventNames,
		},
	)

	for _, ticket := range filteredTickets {
		slog.Info(
			"Found tickets for monitored event",
			"name", ticket.Event.Name,
			"numTickets", ticket.TicketQuantity,
			"ticketPrice", ticket.TotalTicketPrice().String(),
			"originalTicketPrice", ticket.OriginalTicketPrice().String(),
			"link", ticket.Link(),
		)

		err := notificationClient.SendTicketNotification(ticket)
		if err != nil {
			slog.Error(
				"Failed to send notification",
				"err", err,
			)
		}
	}
}
