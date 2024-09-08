package notification_test

import "github.com/ahobsonsayers/twitchets/twickets"

func testNotificationTicket() twickets.Ticket {
	return twickets.Ticket{
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
	}
}
