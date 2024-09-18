package notification

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"text/template"

	"github.com/ahobsonsayers/twitchets/twickets"
)

var (
	//go:embed template/message.tmpl.md
	messageTemplateFS embed.FS
	messageTemplate   *template.Template
)

func init() {
	var err error
	messageTemplate, err = template.ParseFS(messageTemplateFS, "message.tmpl.md")
	if err != nil {
		log.Fatalf("failed to read notification message template: %v", err)
	}
}

type Client interface {
	SendTicketNotification(twickets.Ticket) error
}

type MessageTemplateData struct {
	Venue               string
	Location            string
	Date                string
	Time                string
	NumTickets          int
	TotalTicketPrice    string
	TotalPrice          string
	OriginalTicketPrice string
	OriginalTotalPrice  string
	Discount            float64
	Link                string
}

func renderMessage(ticket twickets.Ticket, includeLink bool) (string, error) {
	templateData := MessageTemplateData{
		Venue:               ticket.Event.Venue.Name,
		Location:            ticket.Event.Venue.Location.Name,
		Date:                ticket.Event.Date.Format("Monday 2 January 2006"),
		Time:                ticket.Event.Time.Format("3:04pm"),
		NumTickets:          ticket.TicketQuantity,
		TotalTicketPrice:    ticket.TotalTicketPrice().String(),
		TotalPrice:          ticket.TotalPrice().String(),
		OriginalTicketPrice: ticket.OriginalTicketPrice().String(),
		OriginalTotalPrice:  ticket.OriginalTotalPrice.String(),
		Discount:            ticket.Discount(),
	}
	if includeLink {
		templateData.Link = ticket.Link()
	}

	var buffer bytes.Buffer
	err := messageTemplate.Execute(&buffer, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render notification message template:, %w", err)
	}

	return buffer.String(), nil
}
