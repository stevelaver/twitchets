package twickets

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

var DefaultClient = NewClient(nil)

// FetchTicketsInput defines parameters when fetching tickets.
// Tickets can either be fetched by number or by time period.
// The default is to get a fixed number of tickets.
// If both a number and time period are set, whichever condition
// is met first will cause tickets fetching to stop.
type FetchTicketsInput struct {
	// Country is required
	Country Country

	Regions []Region

	// Number of tickets to fetch.
	// If getting tickets within in a time period using `CreatedAfter`, set this to an arbitrarily
	// large number (e.g. 250) to ensure all tickets in the period are fetched, while preventing
	// fetching too many tickets and possibly being rate limited or blocked.
	// Defaults to 10.
	// Set to -1 if no limit is desired. This is dangerous and should only be used with well constrained time periods
	NumTickets int

	// Time which fetched tickets must have been created after.
	// Set this to fetch tickets in a time period.
	// Set `NumTickets` to an arbitrarily large number (e.g. 250) to ensure all tickets in the
	// period are fetched, while preventing fetching too many tickets and possibly being rate limited or blocked.
	CreatedAfter time.Time

	// Time which fetched tickets must have been created before.
	// Defaults to current time.
	CreatedBefore time.Time

	// Number of tickets to fetch in each request.
	// Not all tickets are fetched at once - instead
	// a series of requests are made each fetching the
	// number of tickets specified here. In theory this
	// can be arbitrarily large to prevent having to make
	// too many request, however it has been known that
	// any other number than 10 can sometimes not work.
	// Defaults to 10 . Usually can be ignored.
	NumTicketsPerRequest int
}

func (f *FetchTicketsInput) applyDefaults() {
	if f.NumTickets == 0 {
		f.NumTickets = 10
	}
	if f.CreatedBefore.IsZero() {
		f.CreatedBefore = time.Now()
	}
	if f.NumTicketsPerRequest <= 0 {
		f.NumTicketsPerRequest = 10
	}
}

func (f FetchTicketsInput) validate() error {
	if f.Country.Value == "" {
		return errors.New("country must be set")
	}
	if f.CreatedBefore.Before(f.CreatedAfter) {
		return errors.New("created after time must be after the created before time")
	}
	if f.NumTickets < 0 && f.CreatedAfter.IsZero() {
		return errors.New("if not limiting number of tickets, created after must be set")
	}
	return nil
}

// FetchTickets gets the tickets desired by the input struct
func (c *Client) FetchTickets(ctx context.Context, input FetchTicketsInput) (Tickets, error) {
	input.applyDefaults()
	err := input.validate()
	if err != nil {
		return nil, err
	}

	// Iterate through feeds until have equal to or more tickets than desired
	tickets := make(Tickets, 0, input.NumTickets)
	earliestTicketTime := input.CreatedBefore
	for (input.NumTickets < 0 || len(tickets) < input.NumTickets) &&
		earliestTicketTime.After(input.CreatedAfter) {

		feedUrl := FeedUrl(FeedUrlParams{
			Country:    input.Country,
			Regions:    input.Regions,
			NumTickets: input.NumTicketsPerRequest,
			BeforeTime: earliestTicketTime,
		})

		feedTickets, err := c.FetchFeedTickets(ctx, feedUrl)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, feedTickets...)
		earliestTicketTime = feedTickets[len(feedTickets)-1].CreatedAt.Time
	}

	// Only return number of tickets requested
	if len(tickets) > input.NumTickets {
		tickets = tickets[:input.NumTickets]
	}

	// Only return tickets created after the requested time
	if !input.CreatedAfter.IsZero() {
		filteredTickets := make(Tickets, 0, len(tickets))
		for _, ticket := range tickets {
			if ticket.CreatedAt.Time.After(input.CreatedAfter) {
				filteredTickets = append(filteredTickets, ticket)
			}
		}
		tickets = filteredTickets
	}

	return tickets, nil
}

// FetchFeedTickets gets the tickets in the feed url provided.
func (c *Client) FetchFeedTickets(ctx context.Context, feedUrl string) (Tickets, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "") // Twickets blocks some user agents

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		err := fmt.Errorf("error response %s", response.Status)
		if response.StatusCode == http.StatusForbidden {
			err = fmt.Errorf("%s: possibly due to tls misconfiguration", err)
		}
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return UnmarshalTwicketsFeedJson(bodyBytes)
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if httpClient.Transport == nil {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	return &Client{client: httpClient}
}
