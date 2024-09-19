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

	require.Equal(t, twickets.CountryUnitedKingdom, config.Country)

	require.Len(t, config.Events, 4)

	// Event with only name set
	event1 := config.Events[0]
	require.Equal(t, "Event 1", event1.Name)
	require.Equal(t, 0, event1.NumTickets)
	require.Empty(t, event1.Regions)
	require.Equal(t, event1.Discount, 0.0) // nolint: testifylint

	// Event with regions set
	event2 := config.Events[1]
	require.Equal(t, "Event 2", event2.Name)
	require.Len(t, event2.Regions, 1)
	require.Equal(t, event2.Regions[0], twickets.RegionLondon)
	require.Equal(t, 0, event2.NumTickets)
	require.Equal(t, event2.Discount, 0.0) // nolint: testifylint

	// Event with num tickets set
	event3 := config.Events[2]
	require.Equal(t, "Event 3", event3.Name)
	require.Empty(t, event3.Regions)
	require.Equal(t, 1, event3.NumTickets)
	require.Equal(t, event3.Discount, 0.0) // nolint: testifylint

	// Event with discount set
	event4 := config.Events[3]
	require.Equal(t, "Event 4", event4.Name)
	require.Equal(t, 0, event4.NumTickets)
	require.Empty(t, event4.Regions)
	require.Equal(t, event4.Discount, 15.0) // nolint: testifylint
}
