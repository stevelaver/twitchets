package notification_test

import (
	"os"
	"testing"

	"github.com/ahobsonsayers/twitchets/test/testutils"
	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/ahobsonsayers/twitchets/twickets/notification"
	"github.com/stretchr/testify/require"
)

func testNotificationTicket() twickets.Ticket {
	return twickets.Ticket{
		Id: "test",
		Event: twickets.Event{
			Name: "Test Event",
			Venue: twickets.Venue{
				Name: "Test Venue",
				Location: twickets.Location{
					Name: "Test Location",
				},
			},
		},
		TicketType:     "Standing",
		TicketQuantity: 2,
		TicketsPrice: twickets.Price{
			Currency: twickets.CurrencyGBP,
			Amount:   200,
		},
		TwicketsFee: twickets.Price{
			Currency: twickets.CurrencyGBP,
			Amount:   100,
		},
		OriginalTotalPrice: twickets.Price{
			Currency: twickets.CurrencyGBP,
			Amount:   400,
		},
	}
}

func TestRenderMessage(t *testing.T) {
	expectedMessagePath := testutils.ProjectDirectoryJoin(t, "test", "assets", "message.md")
	expectedMessageBytes, err := os.ReadFile(expectedMessagePath)
	require.NoError(t, err)
	expectedMessage := string(expectedMessageBytes)

	tickets := testNotificationTicket()
	actualMessage, err := notification.RenderMessage(tickets)
	require.NoError(t, err)

	require.Equal(t, expectedMessage, actualMessage)
}
