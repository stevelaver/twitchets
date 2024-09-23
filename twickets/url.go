package twickets

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ahobsonsayers/twitchets/twickets/utils"
)

const TwicketsURL = "https://www.twickets.live"

var twicketsUrl *url.URL // https://www.twickets.live

func init() {
	// Parse twickets url
	var err error
	twicketsUrl, err = url.Parse(TwicketsURL)
	if err != nil {
		log.Fatal("failed to parse twickets url")
	}
}

// TicketURL gets a ticket url
// https://www.twickets.live/app/block/<ticketId>,<numTickets>
func TicketURL(ticketId string, numTickets int) string {
	ticketUrl := utils.CloneURL(twicketsUrl)
	ticketUrl = ticketUrl.JoinPath("app", "block", fmt.Sprintf("%s,%d", ticketId, numTickets))
	return ticketUrl.String()
}

type FeedUrlInput struct {
	// Required fields
	APIKey  string
	Country Country

	// Optional fields
	Regions    []Region  // Defaults to all country regions
	NumTickets int       // Defaults to 10 tickets
	BeforeTime time.Time // Defaults to current time
}

func (f FeedUrlInput) validate() error {
	if f.APIKey == "" {
		return errors.New("api key must be set")
	}
	if f.Country.Value == "" {
		return errors.New("country must be set")
	}
	if !Countries.Contains(f.Country) {
		return fmt.Errorf("country '%s' is not valid", f.Country)
	}
	return nil
}

// FeedUrl gets the url of a feed of tickets
// E.g. https://www.twickets.live/services/catalogue?q=countryCode=GB&count=100&api_key=<api_key>
func FeedUrl(input FeedUrlInput) (string, error) {
	err := input.validate()
	if err != nil {
		return "", fmt.Errorf("invalid input parameters: %w", err)
	}

	feedUrl := utils.CloneURL(twicketsUrl)
	feedUrl = feedUrl.JoinPath("services", "catalogue")

	// Set query params
	queryParams := feedUrl.Query()

	locationQuery := apiLocationQuery(input.Country, input.Regions...)
	if locationQuery != "" {
		queryParams.Set("q", locationQuery)
	}

	if !input.BeforeTime.IsZero() {
		maxTime := input.BeforeTime.UnixMilli()
		queryParams.Set("maxTime", strconv.Itoa(int(maxTime)))
	}

	if input.NumTickets > 0 {
		count := strconv.Itoa(input.NumTickets)
		queryParams.Set("count", count)
	}

	queryParams.Set("api_key", input.APIKey)

	// Set query
	encodedQuery := queryParams.Encode()
	encodedQuery = strings.ReplaceAll(encodedQuery, "%3D", "=")
	encodedQuery = strings.ReplaceAll(encodedQuery, "%2C", ",")
	feedUrl.RawQuery = encodedQuery

	return feedUrl.String(), nil
}

// apiLocationQuery converts a country and selection of regions to an api query string
func apiLocationQuery(country Country, regions ...Region) string {
	if !Countries.Contains(country) {
		return ""
	}

	queryParts := make([]string, 0, len(regions)+1)

	countryQuery := fmt.Sprintf("%s=%s", countryQueryKey, country.Value)
	queryParts = append(queryParts, countryQuery)

	for _, region := range regions {
		if Regions.Contains(region) {
			regionQuery := fmt.Sprintf("%s=%s", regionQueryKey, region.Value)
			queryParts = append(queryParts, regionQuery)
		}
	}

	return strings.Join(queryParts, ",")
}
