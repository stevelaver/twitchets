package twickets_test

import (
	"testing"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/stretchr/testify/require"
)

// TestFilterName tests we find tickets for event names that are
// not specified quite right
func TestFilterName(t *testing.T) {
	strangerThingsAsked := "Stranger Things"
	strangerThingsGot := "Stranger Things: The First Shadow"

	backToTheFutureAsked := "Back To The Future"
	backToTheFutureGot := "Back To The Future: The Musical"

	harryPotterAsked := "Harry Potter and the Cursed Child"
	harryPotterGot := "Harry Potter & The Cursed Child Parts 1 & 2"

	wizardOfOzAsked := "The Who"
	wizardOfOzGot := "The The" // This shouldn't match

	gotTickets := twickets.Tickets{
		{Event: twickets.Event{Name: strangerThingsGot}},
		{Event: twickets.Event{Name: backToTheFutureGot}},
		{Event: twickets.Event{Name: harryPotterGot}},
		{Event: twickets.Event{Name: wizardOfOzGot}},
	}

	// Stranger Things
	filteredTickets := gotTickets.Filter([]twickets.Filter{{
		Name: strangerThingsAsked,
	}})
	require.Len(t, filteredTickets, 1)
	require.Equal(t, strangerThingsGot, filteredTickets[0].Event.Name)

	// Back to the Future
	filteredTickets = gotTickets.Filter([]twickets.Filter{{
		Name: backToTheFutureAsked,
	}})
	require.Len(t, filteredTickets, 1)
	require.Equal(t, backToTheFutureGot, filteredTickets[0].Event.Name)

	// Harry Potter
	filteredTickets = gotTickets.Filter([]twickets.Filter{{
		Name: harryPotterAsked,
	}})
	require.Len(t, filteredTickets, 1)
	require.Equal(t, harryPotterGot, filteredTickets[0].Event.Name)

	// Wizard of Oz
	filteredTickets = gotTickets.Filter([]twickets.Filter{{
		Name: wizardOfOzAsked,
	}})
	require.Empty(t, filteredTickets)
}
