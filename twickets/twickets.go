package twickets

import (
	"context"
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

	// Total number of tickets to fetch.
	// Defaults to 100 - large numbers
	// could lead to being rate limited.
	TotalNumTickets int

	// Number of tickets in fetch in each feed.
	// Defaults to 10 - other numbers have been known to fail
	NumTicketsPerFeed int

	// Time to get tickets before.
	// Defaults to current time.
	BeforeTime time.Time
}

func (f *FetchTicketsInput) applyDefaults() {
	if f.TotalNumTickets <= 0 {
		f.TotalNumTickets = 100
	}
	if f.NumTicketsPerFeed <= 0 {
		f.NumTicketsPerFeed = 10
	}
	if f.BeforeTime.IsZero() {
		f.BeforeTime = time.Now()
	}
}

// FetchTickets gets the tickets desired by the input struct
func (c *Client) FetchTickets(ctx context.Context, input FetchTicketsInput) (Tickets, error) {
	input.applyDefaults()

	// Iterate through feeds until have equal to or more tickets than desired
	tickets := make(Tickets, 0, input.TotalNumTickets)
	earliestTime := input.BeforeTime
	for len(tickets) < input.TotalNumTickets {

		feedUrl := FeedUrl(FeedUrlParams{
			Country:    input.Country,
			Regions:    input.Regions,
			NumTickets: input.NumTicketsPerFeed,
			BeforeTime: earliestTime,
		})

		feedTickets, err := c.FetchFeedTickets(ctx, feedUrl)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, feedTickets...)
		earliestTime = feedTickets[len(feedTickets)-1].CreatedAt.Time
	}

	// Return the exact number of tickets asked for
	return tickets[:input.TotalNumTickets], nil
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
		return nil, fmt.Errorf("got error response: %s", response.Status)
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
	return &Client{client: httpClient}
}
