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
	Country    Country
	Regions    []Region
	MaxNumber  int
	BeforeTime time.Time
}

func DefaultFetchTicketsInput(country Country) FetchTicketsInput {
	return FetchTicketsInput{
		Country:   country,
		MaxNumber: 10,
	}
}

// FetchLatestTickets gets latest listed tickets in a country and region(s) since the
// up to a maximum limit, before any specified time (defaults to now)
func (c *Client) FetchLatestTickets(ctx context.Context, input FetchTicketsInput) ([]Ticket, error) {
	feedUrl := FeedUrl(input)
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
