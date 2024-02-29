package response

type Ticket struct {
	TicketType    string `json:"ticketType"`
	TicketPrice   string `json:"ticketPrice"`
	ContinentName string `json:"continentName"`
	ContinentCode string `json:"continentCode"`
	CountryName   string `json:"countryName"`
	CountryCode   string `json:"countryCode"`
	IsSold        bool   `json:"isSold"`
}

type SuggestionTicket struct {
	TicketType          string `json:"ticketType"`
	NormalTicketPrice   string `json:"normalTicketPrice"`
	DiscountTicketPrice string `json:"discountTicketPrice"`
	Discount            string `json:"discount"`
	ContinentName       string `json:"continentName"`
	ContinentCode       string `json:"continentCode"`
	CountryName         string `json:"countryName"`
	CountryCode         string `json:"countryCode"`
	IsSold              bool   `json:"isSold"`
}

type TicketResp struct {
	Tickets    []Ticket           `json:"tickets"`
	Suggestion []SuggestionTicket `json:"suggestion"`
}

type TicketCountry struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	IsSold      bool   `json:"isSold"`
}
