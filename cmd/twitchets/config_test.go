package main

import (
	"testing"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	configPath := testutils.ProjectDirectoryJoin(t, "test", "assets", "config", "config.yaml")
	actualConfig, err := LoadConfig(configPath)
	require.NoError(t, err)

	globalCountry := twickets.CountryUnitedKingdom
	globalRegions := []twickets.Region{twickets.RegionLondon, twickets.RegionNorthWest}
	globalNumTickets := 2
	globalDiscount := 25.0

	expectedConfig := Config{
		APIKey: "test",
		GlobalConfig: GlobalEventConfig{
			Country:    globalCountry,
			Regions:    globalRegions,
			NumTickets: globalNumTickets,
			Discount:   globalDiscount,
		},
		TicketsConfig: []TicketConfig{
			{
				// Event with only name set
				Name: "Event 1",
			},
			{
				// Event with regions set
				Name:    "Event 2",
				Regions: []twickets.Region{twickets.RegionSouthWest},
			},
			{
				// Event with num tickets set
				Name:       "Event 3",
				NumTickets: lo.ToPtr(1),
			},
			{
				// Event with discount set
				Name:     "Event 4",
				Discount: lo.ToPtr(15.0),
			},
		},
	}

	require.EqualValues(t, expectedConfig, actualConfig)
}

func TestConfigFilters(t *testing.T) {
	configPath := testutils.ProjectDirectoryJoin(t, "test", "assets", "config", "config.yaml")
	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	actualFilters := config.Filters()

	globalRegions := []twickets.Region{twickets.RegionLondon, twickets.RegionNorthWest}
	globalNumTickets := 2
	globalDiscount := 25.0

	expectedFilters := []twickets.Filter{
		{
			// Event with only name set
			Name:       "Event 1",
			Regions:    globalRegions,
			NumTickets: globalNumTickets,
			Discount:   globalDiscount,
		},
		{
			// Event with regions set
			Name:       "Event 2",
			Regions:    []twickets.Region{twickets.RegionSouthWest},
			NumTickets: globalNumTickets,
			Discount:   globalDiscount,
		},
		{
			// Event with num tickets set
			Name:       "Event 3",
			Regions:    globalRegions,
			NumTickets: 1,
			Discount:   globalDiscount,
		},
		{
			// Event with discount set
			Name:       "Event 4",
			Regions:    globalRegions,
			NumTickets: globalNumTickets,
			Discount:   15.0,
		},
		{
			// Event with globals unset
			Name:       "Event 5",
			Regions:    []twickets.Region{},
			NumTickets: 0,
			Discount:   0.0,
		},
	}

	require.EqualValues(t, expectedFilters, actualFilters)
}
