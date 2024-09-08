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
	lines := []string{
		fmt.Sprintf(
			"%s %s",
			ticket.Event.Date.Format("Monday 2 January 2006"),
			ticket.Event.Time.Format("3:04pm"),
		),
		fmt.Sprintf("%d ticket(s)", ticket.TicketQuantity),
		"",
		fmt.Sprintf("Ticket Price: %s", ticket.TotalTicketPrice().String()),
		fmt.Sprintf("Total Price: %s", ticket.TotalPrice().String()),
		fmt.Sprintf("Discount: %s", ticket.DiscountString()),
		"",
		fmt.Sprintf("Original Ticket Price: %s", ticket.OriginalTicketPrice().String()),
		fmt.Sprintf("Original Total Price: %s", ticket.OriginalTotalPrice.String()),
	}
	if includeLink {
		lines = append(lines,
			"",
			fmt.Sprintf("[Buy](%s)", ticket.Link()),
		)
	}

	return strings.Join(lines, "\n")
}
