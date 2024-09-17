package twickets_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

func TestGetLatestTickets(t *testing.T) {
	testutils.SkipIfCI(t, "Does not work in CI. Fix")

	httpClient := http.DefaultClient
	// httpClient, err := testutils.NewProxyClient(testutils.RoosterKidProxyListURL)
	// require.NoError(t, err)

	twicketsClient := twickets.NewClient(httpClient)
	tickets, err := twicketsClient.FetchTickets(
		context.Background(),
		twickets.FetchTicketsInput{
			Country: twickets.CountryUnitedKingdom,
			Regions: []twickets.Region{
				twickets.RegionLondon,
				twickets.RegionNorthWest,
			},
			CreatedAfter: time.Now().Add(-10 * time.Minute),
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, tickets)
}
