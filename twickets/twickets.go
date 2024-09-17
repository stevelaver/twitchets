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

type FetchTicketsInput struct {
	Country Country
	Regions []Region

	// Time which tickets must have been created before.
	// Defaults to current time.
	CreatedBefore time.Time

	// Time which tickets must have been created after.
	// Defaults to a minute before the current time.
	CreatedAfter time.Time

	// Max number of tickets to fetch in the time period.
	// Set this to an arbitrarily large number to prevent
	// fetching too many tickets at once and possibly
	// being rate limited.
	// Defaults to 250. Usually can be ignored.
	MaxTickets int

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
	if f.CreatedBefore.IsZero() {
		f.CreatedBefore = time.Now()
	}
	if f.CreatedAfter.IsZero() {
		f.CreatedAfter = f.CreatedBefore.Add(-time.Minute)
	}
	if f.MaxTickets <= 0 {
		f.MaxTickets = 250
	}
	if f.NumTicketsPerRequest <= 0 {
		f.NumTicketsPerRequest = 10
	}
}

// FetchTickets gets the tickets desired by the input struct
func (c *Client) FetchTickets(ctx context.Context, input FetchTicketsInput) (Tickets, error) {
	input.applyDefaults()

	if input.Country.Value == "" {
		return nil, errors.New("country must be set")
	}
	if input.CreatedBefore.Before(input.CreatedAfter) {
		return nil, errors.New("latests time must be after earliest time")
	}

	// Iterate through feeds until have equal to or more tickets than desired
	tickets := make(Tickets, 0, input.MaxTickets)
	earliestTicketTime := input.CreatedBefore
	for earliestTicketTime.After(input.CreatedAfter) && len(tickets) < input.MaxTickets {

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

	// Only return tickets created after the earliest time
	tickets = tickets.ticketsCreatedAfterTime(input.CreatedAfter)

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
