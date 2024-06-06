package twickets // nolint

import (
	"encoding/json"
)

type Ticket struct {
	// Buy Link: https://www.twickets.live/app/block/<ticketId>,<quanitity>
	Id string `json:"blockId"`

	CreatedAt UnixTime `json:"created"`
	ExpiresAt UnixTime `json:"expires"`

	TicketQuantity           int   `json:"ticketQuantity"`
	TotalSellingPrice        Price `json:"totalSellingPrice"`
	TotalTwicketsFee         Price `json:"totalTwicketsFee"`
	FaceValuePrice           Price `json:"faceValuePrice"`
	SellerWillConsiderOffers bool  `json:"sellerWillConsiderOffers"`

	TicketType   string `json:"priceTier"` // Seated, Standing, Box etc.
	SeatAssigned bool   `json:"seatAssigned"`
	Section      string `json:"section"` // Can be empty
	Row          string `json:"row"`     // Can be empty

	Event Event `json:"event"`
	Tour  Tour  `json:"tour"`
}

type Event struct {
	ID       string `json:"id"`
	Name     string `json:"eventName"`
	Category string `json:"category"`
	Date     Date   `json:"date"`
	Time     Time   `json:"showStartingTime"`

	OnSaleTime    *DateTime `json:"onSaleTime"` // 2023-11-17T10:00:00Z
	AnnouncedTime *DateTime `json:"created"`    // 2023-11-17T10:00:00Z

	Venue  Venue    `json:"venue"`
	Lineup []Lineup `json:"participants"`
}

type Lineup struct {
	Artist  Artist `json:"participant"`
	Billing int    `json:"billing"`
}

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"linkName"`
}

type Venue struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Postcode string   `json:"postcode"`
}

type Location struct {
	Id          string `json:"id"`
	Name        string `json:"shortName"`
	FullName    string `json:"name"`
	CountryCode string `json:"countryCode"` // TODO use enum
	RegionCode  string `json:"regionCode"`  // TODO use enum
}

type Tour struct {
	Id           string   `json:"tourId"`
	Name         string   `json:"tourName"`
	Slug         string   `json:"slug"`
	FirstEvent   *Date    `json:"minDate"`      // 2024-06-06
	LastEvent    *Date    `json:"maxDate"`      // 2024-11-14"
	CountryCodes []string `json:"countryCodes"` // TODO use enum
}

func UnmarshalTwicketsFeedJson(data []byte) ([]Ticket, error) {
	response := struct {
		ResponseData []struct { // nolint
			Tickets *Ticket `json:"catalogBlockSummary"`
		} `json:"responseData"`
	}{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	// Get non null tickets. Tickets are null if they have been delisted
	tickets := make([]Ticket, 0, len(response.ResponseData))
	for _, responseData := range response.ResponseData {
		if responseData.Tickets != nil {
			tickets = append(tickets, *responseData.Tickets)
		}
	}

	return tickets, nil
}
