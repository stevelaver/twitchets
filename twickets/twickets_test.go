package twickets_test

import (
	"context"
	"net/http"
	"testing"

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
			NumTickets: 10,
		},
	)
	require.NoError(t, err)
	require.Len(t, tickets, 10)
}
