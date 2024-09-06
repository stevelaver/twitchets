package twickets_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

func TestGetLatestTickets(t *testing.T) {
	// t.Skip("Does not work in CI")

	httpClient, err := testutils.NewProxyClient(testutils.RoosterKidProxyListURL)
	require.NoError(t, err)

	twicketsClient := twickets.NewClient(httpClient)
	tickets, err := twicketsClient.FetchLatestTickets(
		context.Background(),
		twickets.DefaultFetchTicketsInput(twickets.CountryUnitedKingdom),
	)
	require.NoError(t, err)
	require.Len(t, tickets, 10)
}
