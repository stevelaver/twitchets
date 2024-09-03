package twickets

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

const (
	GotifyURL   = "https://notifications.arranhs.com"
	GotifyToken = "AxNVvRfx9.ZKCTj"
)

var (
	gotifyURL    *url.URL
	gotifyClient *client.GotifyREST
)

func init() {
	// Parse gotify url
	var err error
	gotifyURL, err = url.Parse(GotifyURL)
	if err != nil {
		log.Fatal("failed to parse gotify url")
	}

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
		auth.TokenAuth(GotifyToken),
	)
	if err != nil {
		return err
	}

	return nil
}
