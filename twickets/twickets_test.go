package twickets_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

func TestGetLatestTickets(t *testing.T) {
	t.Skip("Does not work in CI")

	// client, err := testutils.NewProxyClient(testutils.ProxyListURL)
	// require.NoError(t, err)

	client := twickets.NewClient(nil)
	tickets, err := client.FetchLatestTickets(
		context.Background(),
		twickets.DefaultFetchTicketsInput(twickets.CountryUnitedKingdom),
	)
	require.NoError(t, err)
	require.Len(t, tickets, 10)
}
