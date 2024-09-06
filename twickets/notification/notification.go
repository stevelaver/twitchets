package notification

import (
	"fmt"
	"strings"

	"github.com/ahobsonsayers/twitchets/twickets"
)

type Client interface {
	SendTicketNotification(twickets.Ticket) error
}

func notificationMessage(ticket twickets.Ticket) string {
	lines := []string{
		fmt.Sprintf(
			"%s %s",
			ticket.Event.Time.Format("3:04pm"),
			ticket.Event.Date.Format("Monday 2 January 2006"),
		),
		fmt.Sprintf("%d ticket(s)", ticket.TicketQuantity),
		"",
		fmt.Sprintf("Ticket Price: %s", ticket.TotalTicketPrice().String()),
		fmt.Sprintf("Original Ticket Price: %s", ticket.OriginalTicketPrice().String()),
		"",
		fmt.Sprintf("Total Price: %s", ticket.TotalPrice().String()),
		fmt.Sprintf("Original Total Price: %s", ticket.OriginalTotalPrice.String()),
		"",
		fmt.Sprintf("[Buy](%s)", ticket.Link()),
	}

	return strings.Join(lines, "\n")
}
