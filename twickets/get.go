package twickets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetTicketsInput struct {
	Country    Country
	Regions    []Region
	MaxNumber  int
	BeforeTime time.Time
}

// GetLatestTickets gets latest listed tickets in a country and region(s) since the
// up to a maximum limit, before any specified time (defaults to now)
func GetLatestTickets(ctx context.Context, input GetTicketsInput) ([]Ticket, error) {
	feedUrl := FeedUrl(input)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response != nil && response.StatusCode >= 300 {
		return nil, fmt.Errorf("got error response: %s", response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return UnmarshalTwicketsFeedJson(bodyBytes)
}
