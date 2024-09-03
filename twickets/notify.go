package twickets

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/joho/godotenv"
)

var (
	gotifyURL    *url.URL
	gotifyToken  string
	gotifyClient *client.GotifyREST
)

func init() {
	// Load .env file if it exists
	_ = godotenv.Load()

	gotifyUrlEnvVar := os.Getenv("GOTIFY_URL")
	if gotifyUrlEnvVar == "" {
		log.Fatal("GOTIFY_URL is not set")
	}

	gotifyTokenEnvVar := os.Getenv("GOTIFY_TOKEN")
	if gotifyTokenEnvVar == "" {
		log.Fatal("GOTIFY_TOKEN is not set")
	}

	var err error
	gotifyURL, err = url.Parse(gotifyUrlEnvVar)
	if err != nil {
		log.Fatal("failed to parse gotify url")
	}

	gotifyToken = gotifyTokenEnvVar // TODO validate this token

	gotifyClient = gotify.NewClient(gotifyURL, &http.Client{})
}

func SendTicketNotification(ticket Ticket) error {
	notificationMessage := fmt.Sprintf(`
	Day: %s
	Number of tickets: %d
	Ticket Price: %s
	Orignal Price: %s
	`,
		ticket.Event.Date.Format("Monday 2 January 2006"),
		ticket.TicketQuantity,
		ticket.TotalSellingPrice.PerString(ticket.TicketQuantity),
		ticket.TotalSellingPrice.String(),
	)

	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    ticket.Event.Name,
		Message:  notificationMessage,
		Priority: 5,
	}

	_, err := gotifyClient.Message.CreateMessage(
		params,
		auth.TokenAuth(gotifyToken),
	)
	if err != nil {
		return err
	}

	return nil
}
