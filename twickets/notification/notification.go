package notification

import (
	"fmt"
	"strings"

	"github.com/ahobsonsayers/twitchets/twickets"
)

type Client interface {
	SendTicketNotification(twickets.Ticket) error
}

func notificationMessage(ticket twickets.Ticket, includeLink bool) string { // nolint: revive
	var builder strings.Builder

	writeLine(&builder, "%s, %s", ticket.Event.Venue.Name, ticket.Event.Venue.Location.Name)
	writeLine(&builder, "%s %s", ticket.Event.Date.Format("Monday 2 January 2006"), ticket.Event.Time.Format("3:04pm"))
	writeLine(&builder, "%d ticket(s)", ticket.TicketQuantity)

	writeLine(&builder, "")

	writeLine(&builder, "Ticket Price: %s", ticket.TotalTicketPrice().String())
	writeLine(&builder, "Total Price: %s", ticket.TotalPrice().String())
	if ticket.Discount() < 0 {
		writeLine(&builder, "Discount: None")
	} else {
		writeLine(&builder, "Discount: %s", ticket.DiscountString())
	}

	writeLine(&builder, "")

	writeLine(&builder, "Original Ticket Price: %s", ticket.OriginalTicketPrice().String())
	writeLine(&builder, "Original Total Price: %s", ticket.OriginalTotalPrice.String())

	if includeLink {
		writeLine(&builder, "")
		writeLine(&builder, "Buy: %s", ticket.Link())
	}

	return builder.String()
}

func writeLine(builder *strings.Builder, format string, args ...any) {
	_, _ = builder.WriteString(
		fmt.Sprintf(format, args...) + "\n",
	)
}
