package main

import (
	"path/filepath"
	"testing"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

func getTestConfigDirectory(t *testing.T) string {
	projectDirectory := testutils.GetProjectDirectory(t)
	return filepath.Join(projectDirectory, "test", "assets", "config")
}

func TestLoadConfig(t *testing.T) {
	configPath := filepath.Join(getTestConfigDirectory(t), "config.yaml")
	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	globalCountry := twickets.CountryUnitedKingdom
	globalRegions := []twickets.Region{twickets.RegionLondon, twickets.RegionNorthWest}
	globalNumTickets := 2
	globalDiscount := 25

	require.Equal(t, globalCountry, config.GlobalConfig.Country)
	require.Equal(t, globalRegions, config.GlobalConfig.Regions)
	require.Equal(t, globalNumTickets, config.GlobalConfig.NumTickets)
	require.InDelta(t, globalDiscount, config.GlobalConfig.Discount, 0)

	require.Len(t, config.EventConfig, 4)

	// Event with only name set
	event1 := config.EventConfig[0]
	// Global Config
	require.Equal(t, globalRegions, event1.Regions)
	require.Equal(t, globalNumTickets, event1.NumTickets)
	require.InDelta(t, globalDiscount, event1.Discount, 0)
	// Event config
	require.Equal(t, "Event 1", event1.Name)

	// Event with regions set
	event2 := config.EventConfig[1]
	// Global Config
	require.Equal(t, globalNumTickets, event2.NumTickets)
	require.InDelta(t, globalDiscount, event2.Discount, 0)
	// Event config
	require.Equal(t, "Event 2", event2.Name)
	require.Len(t, event2.Regions, 1)
	require.Equal(t, event2.Regions[0], twickets.RegionSouthWest)

	// Event with num tickets set
	event3 := config.EventConfig[2]
	// Global Config
	require.Equal(t, globalRegions, event3.Regions)
	require.InDelta(t, globalDiscount, event3.Discount, 0)
	// Event config
	require.Equal(t, "Event 3", event3.Name)
	require.Equal(t, 1, event3.NumTickets)

	// Event with discount set
	event4 := config.EventConfig[3]
	// Global Config
	require.Equal(t, globalRegions, event4.Regions)
	require.Equal(t, globalNumTickets, event4.NumTickets)
	// Event config
	require.Equal(t, "Event 4", event4.Name)
	require.InDelta(t, 15.0, event4.Discount, 0) // nolint: testifylint
}
