package request

type TicketReq struct {
	CountryCode string `json:"countryCode" validate:"required"`
	EventId     string `json:"eventId" validate:"required"`
}

type CreateOnlineTicketReq struct {
	Tag         string `json:"tag" validate:"required"`
	CountryCode string `json:"countryCode" validate:"required"`
}
