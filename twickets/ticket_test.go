package twickets_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalFeedJson(t *testing.T) {
	projectDirectory := testutils.GetProjectDirectory(t)
	feedJsonFilePath := filepath.Join(projectDirectory, "test", "assets", "feed.json")

	feedJsonFile, err := os.Open(feedJsonFilePath)
	require.NoError(t, err)
	feedJson, err := io.ReadAll(feedJsonFile)
	require.NoError(t, err)

	tickets, err := twickets.UnmarshalTwicketsFeedJson(feedJson)
	require.NoError(t, err)

	require.Len(t, tickets, 4)

	require.Equal(t, tickets[0].Event.Name, "Foo Fighters")
	require.Len(t, tickets[0].Event.Lineup, 3)
	require.Equal(t, tickets[0].Event.Lineup[0].Artist.Name, "Foo Fighters")
	require.Equal(t, tickets[0].Event.Lineup[1].Artist.Name, "Wet Leg")
	require.Equal(t, tickets[0].Event.Lineup[2].Artist.Name, "Shame")
	require.Equal(t, tickets[0].Event.Venue.Name, "London Stadium")
	require.Equal(t, tickets[0].Tour.Name, "Foo Fighters - Everything Or Nothing At All Tour")
	require.Equal(t, tickets[0].TicketQuantity, 3)
	require.Equal(t, tickets[0].TotalSellingPrice.String(), "£255.00")
	require.Equal(t, tickets[0].TotalTwicketsFee.String(), "£38.25")
	require.Equal(t, tickets[0].FaceValuePrice.String(), "£255.00")

	require.Equal(t, tickets[1].Event.Name, "Mean Girls")
	require.Len(t, tickets[1].Event.Lineup, 0)
	require.Equal(t, tickets[1].Event.Venue.Name, "Savoy Theatre")
	require.Equal(t, tickets[1].Tour.Name, "Mean Girls")
	require.Equal(t, tickets[1].TicketQuantity, 2)
	require.Equal(t, tickets[1].TotalSellingPrice.String(), "£130.00")
	require.Equal(t, tickets[1].TotalTwicketsFee.String(), "£18.20")
	require.Equal(t, tickets[1].FaceValuePrice.String(), "£130.00")

	require.Equal(t, tickets[2].Event.Name, "South Africa v Wales")
	require.Len(t, tickets[2].Event.Lineup, 0)
	require.Equal(t, tickets[2].Event.Venue.Name, "Twickenham Stadium")
	require.Equal(t, tickets[2].Tour.Name, "South Africa v Wales")
	require.Equal(t, tickets[2].TicketQuantity, 4)
	require.Equal(t, tickets[2].TotalSellingPrice.String(), "£380.00")
	require.Equal(t, tickets[2].TotalTwicketsFee.String(), "£53.20")
	require.Equal(t, tickets[2].FaceValuePrice.String(), "£380.00")

	require.Equal(t, tickets[3].Event.Name, "Download Festival 2024")
	require.Len(t, tickets[3].Event.Lineup, 0)
	require.Equal(t, tickets[3].Event.Venue.Name, "Donington Park")
	require.Equal(t, tickets[3].Tour.Name, "Download Festival 2024")
	require.Equal(t, tickets[3].TicketQuantity, 1)
	require.Equal(t, tickets[3].TotalSellingPrice.String(), "£280.00")
	require.Equal(t, tickets[3].TotalTwicketsFee.String(), "£30.80")
	require.Equal(t, tickets[3].FaceValuePrice.String(), "£322.00")
}
