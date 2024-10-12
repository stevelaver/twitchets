package twickets_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGetLatestTickets(t *testing.T) {
	testutils.SkipIfCI(t, "This has issues in CI. Proxy client doesnt work properly yet")

	_ = godotenv.Load(testutils.ProjectDirectoryJoin(t, ".env"))

	twicketsAPIKey := os.Getenv("TWICKETS_API_KEY")
	require.NotEmpty(t, twicketsAPIKey, "TWICKETS_API_KEY is not set")

	// Use proxy client in CI
	var httpClient *http.Client
	if !testutils.IsCI() {
		httpClient = http.DefaultClient
	} else {
		var err error
		httpClient, err = testutils.NewProxyClient(
			testutils.ProxyClientConfig{
				ProxyListUrls: testutils.ProxyListUrls,
				TestTimeout:   5 * time.Second,
			},
		)
		require.NoError(t, err)
	}

	twicketsClient := twickets.NewClient(httpClient)
	tickets, err := twicketsClient.FetchTickets(
		context.Background(),
		twickets.FetchTicketsInput{
			APIKey:  twicketsAPIKey,
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
