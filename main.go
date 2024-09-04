package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/joho/godotenv"
)

const maxNumTickets = 10

var (
	// Config variables
	countryCode         = "GB"
	regionCodes         = []string{"GBLO"}
	monitoredEventNames = []string{
		// Theatre
		"Back to the Future",
		"Frozen",
		"Hamilton",
		"Harry Potter & the Cursed Child",
		"Kiss Me Kate",
		"Lion King",
		"Matilda",
		"Mean Girls",
		"Moulin Rouge",
		"Starlight Express",
		"Stranger Things",
		"The Phantom Opera",
		"The Wizard of Oz",
		// Gigs
		"Coldplay",
		"Taylor Swift",
	}

	// Package variables
	lastCheckTime = time.Time{}
	country       twickets.Country
	regions       []twickets.Region
)

func init() {
	_ = godotenv.Load()
}

func main() {
	gotifyUrl := os.Getenv("GOTIFY_URL")
	if gotifyUrl == "" {
		log.Fatal("GOTIFY_URL is not set")
	}

	gotifyToken := os.Getenv("GOTIFY_TOKEN")
	if gotifyToken == "" {
		log.Fatal("GOTIFY_TOKEN is not set")
	}

	country := twickets.Countries.Parse(countryCode)
	if country == nil {
		log.Fatalf("'%s' is not a valid country code", country)
	}

	regions = make([]twickets.Region, 0, len(regionCodes))
	for _, regionCode := range regionCodes {
		region := twickets.Regions.Parse(regionCode)
		if region == nil {
			log.Fatalf("'%s' is not a valid region code", region)
		}
	}

	notificationClient, err := twickets.NewGotifyClient(gotifyUrl, gotifyToken)
	if err != nil {
		log.Fatal(err)
	}

	twicketsClient := twickets.NewClient(nil)

	slog.Info(
		fmt.Sprintf("Monitoring: %s", strings.Join(monitoredEventNames, ", ")),
	)

	// Initial execution
	fetchAndProcessTickets(twicketsClient, notificationClient)

	// Create ticker
	ticker := time.NewTicker(1 * time.Minute)
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
	notificationClient twickets.NotificationClient,
) {
	checkTime := time.Now()
	defer func() { lastCheckTime = checkTime }()

	tickets, err := twicketsClient.FetchLatestTickets(
		context.Background(),
		twickets.FetchTicketsInput{
			Country:    country,
			Regions:    regions,
			MaxNumber:  maxNumTickets,
			BeforeTime: checkTime,
		},
	)
	if err != nil {
		// If there is an error, try again with a no input struct
		// Twickets api has been known to fail with other query params outside the defaults
		tickets, err = twicketsClient.FetchLatestTickets(
			context.Background(),
			twickets.DefaultFetchTicketsInput(country),
		)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	filteredTickets := twickets.FilterTickets(
		tickets,
		twickets.TicketFilter{
			EventNames:   monitoredEventNames,
			CreatedAfter: lastCheckTime,
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
