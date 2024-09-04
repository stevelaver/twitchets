package twickets

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

type NotificationClient interface {
	SendTicketNotification(Ticket) error
}

func notificationMessage(ticket Ticket) string {
	lines := []string{
		ticket.Event.Date.Format("Monday 2 January 2006"),
		fmt.Sprintf("%d ticket(s)", ticket.TicketQuantity),
		fmt.Sprintf("Ticket Price: %s", ticket.TotalSellingPrice.PerString(ticket.TicketQuantity)),
		fmt.Sprintf("Original Price: %s", ticket.FaceValuePrice.PerString(ticket.TicketQuantity)),
	}

	return strings.Join(lines, "\n")
}

type GotifyClient struct {
	url    *url.URL
	token  string
	client *client.GotifyREST
}

var _ NotificationClient = GotifyClient{}

func (g GotifyClient) SendTicketNotification(ticket Ticket) error {
	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:    ticket.Event.Name,
		Message:  notificationMessage(ticket),
		Extras:   map[string]any{},
		Priority: 5,
	}

	_, err := g.client.Message.CreateMessage(
		params,
		auth.TokenAuth(g.token),
	)
	if err != nil {
		return err
	}

	return nil
}

func NewGotifyClient(gotifyUrl, gotifyToken string) (*GotifyClient, error) {
	parsedGotifyUrl, err := url.Parse(gotifyUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gotify url: %v", err)
	}

	return &GotifyClient{
		url:    parsedGotifyUrl,
		token:  gotifyToken, // TODO validate this token?
		client: gotify.NewClient(parsedGotifyUrl, &http.Client{}),
	}, nil
}
