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

	require.Len(t, config.Regions, 2)
	require.Equal(t, twickets.RegionLondon, config.Regions[0])
	require.Equal(t, twickets.RegionNorthWest, config.Regions[1])

	require.Len(t, config.Events, 2)

	require.Equal(t, "Event 1", config.Events[0].Name)
	require.Equal(t, 0, config.Events[0].NumTickets)

	require.Equal(t, "Event 2", config.Events[1].Name)
	require.Equal(t, 1, config.Events[1].NumTickets)
}
