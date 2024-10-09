package notification

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"strings"
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
	messageTemplate, err = template.ParseFS(messageTemplateFS, "template/message.tmpl.md")
	if err != nil {
		log.Fatalf("failed to read notification message template: %v", err)
	}
}

type Client interface {
	SendTicketNotification(twickets.Ticket) error
}

type MessageTemplateData struct {
	// Header
	Event string

	// Body
	Date                string
	Time                string
	Venue               string
	Location            string
	TicketType          string // Standing, Stalls etc.
	NumTickets          int
	TotalTicketPrice    string
	TotalPrice          string
	OriginalTicketPrice string
	OriginalTotalPrice  string
	Discount            float64

	// Footer
	Link string
}

type RenderMessageConfig struct {
	// Whether to include event name
	IncludeHeader bool

	// Whether to include buy link
	IncludeFooter bool
}

func RenderMessage(ticket twickets.Ticket, config *RenderMessageConfig) (string, error) {
	templateData := MessageTemplateData{
		Date:                ticket.Event.Date.Format("Monday 2 January 2006"),
		Time:                ticket.Event.Time.Format("3:04pm"),
		Venue:               ticket.Event.Venue.Name,
		Location:            ticket.Event.Venue.Location.Name,
		TicketType:          ticket.TicketType,
		NumTickets:          ticket.TicketQuantity,
		TotalTicketPrice:    ticket.TotalTicketPrice().String(),
		TotalPrice:          ticket.TotalPrice().String(),
		OriginalTicketPrice: ticket.OriginalTicketPrice().String(),
		OriginalTotalPrice:  ticket.OriginalTotalPrice.String(),
		Discount:            ticket.Discount(),
	}

	// add optional header and footers
	if config != nil {
		if config.IncludeHeader {
			templateData.Event = ticket.Event.Name
		}

		if config.IncludeFooter {
			templateData.Link = ticket.Link()
		}
	}

	var buffer bytes.Buffer
	err := messageTemplate.Execute(&buffer, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render notification message template:, %w", err)
	}

	message := buffer.String()
	message = strings.TrimSpace(message)

	return message, nil
}
