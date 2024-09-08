package twickets

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ahobsonsayers/twitchets/twickets/utils"
)

const (
	TwicketsURL    = "https://www.twickets.live"
	TwicketsAPIKey = "83d6ec0c-54bb-4da3-b2a1-f3cb47b984f1"
)

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

type FeedUrlParams struct {
	Country Country
	Regions []Region

	// Number of tickets to fetch in the feed.
	NumTickets int

	// Time to get tickets before.
	BeforeTime time.Time
}

// FeedUrl gets the url of the feed with the given params.
// E.g. https://www.twickets.live/services/catalogue?q=countryCode=GB&count=100&api_key=<api_key>
func FeedUrl(input FeedUrlParams) string {
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

	queryParams.Set("api_key", TwicketsAPIKey)

	// Set query
	encodedQuery := queryParams.Encode()
	encodedQuery = strings.ReplaceAll(encodedQuery, "%3D", "=")
	encodedQuery = strings.ReplaceAll(encodedQuery, "%2C", ",")
	feedUrl.RawQuery = encodedQuery

	return feedUrl.String()
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
