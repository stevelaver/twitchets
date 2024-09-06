package notification_test

import (
	"os"
	"testing"

	"github.com/ahobsonsayers/twitchets/twickets"
	"github.com/ahobsonsayers/twitchets/twickets/notification"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"heckel.io/ntfy/log"
)

func TestNtfySendTicketMessage(t *testing.T) {
	t.Skip("No env set in CI. Fix")
	log.SetLevel(log.DebugLevel)
	_ = godotenv.Load("../../.env")

	ntfyUrl := os.Getenv("NTFY_URL")
	require.NotEmpty(t, ntfyUrl, "NTFY_URL is not set")

	ntfyUser := os.Getenv("NTFY_USER")
	require.NotEmpty(t, ntfyUser, "NTFY_USER is not set")

	ntfyPassword := os.Getenv("NTFY_PASSWORD")
	require.NotEmpty(t, ntfyPassword, "NTFY_PASSWORD is not set")

	ntfyTopic := os.Getenv("NTFY_TOPIC")
	require.NotEmpty(t, ntfyTopic, "NTFY_PASSWORD is not set")

	client, err := notification.NewNtfyClient(notification.NtfyConfig{
		Url:      ntfyUrl,
		User:     ntfyUser,
		Password: ntfyPassword,
		Topic:    ntfyTopic,
	})
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
