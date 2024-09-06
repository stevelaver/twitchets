package notification_test

import (
	"os"
	"testing"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/ahobsonsayers/twitchets/twickets/notification"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGotifySendTicketMessage(t *testing.T) {
	t.Skip("No env set in CI. Fix")
	_ = godotenv.Load("../../.env")

	gotifyUrl := os.Getenv("GOTIFY_URL")
	require.NotEmpty(t, gotifyUrl, "GOTIFY_URL is not set")

	gotifyToken := os.Getenv("GOTIFY_TOKEN")
	require.NotEmpty(t, gotifyToken, "GOTIFY_TOKEN is not set")

	client, err := notification.NewGotifyClient(gotifyUrl, gotifyToken)
	require.NoError(t, err)

	err = client.SendTicketNotification(twickets.Ticket{
		Id: "test",
		Event: twickets.Event{
			Name: "Test Event",
		},
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
	})
	require.NoError(t, err)
}
