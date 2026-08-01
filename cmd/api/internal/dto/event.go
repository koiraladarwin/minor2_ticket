package dto

import "time"

type CreateEventRequest struct {
	Title             string    `json:"title"`
	Description       *string   `json:"description,omitempty"`
	Venue             string    `json:"venue"`
	BannerURL         *string   `json:"banner_url,omitempty"`
	EventStartAt      time.Time `json:"event_start_at"`
	EventEndAt        time.Time `json:"event_end_at"`
	TicketSaleStartAt time.Time `json:"ticket_sale_start_at"`
	TicketSaleEndAt   time.Time `json:"ticket_sale_end_at"`
	Capacity          int       `json:"capacity"`
}

type UpdateEventRequest struct {
	Title             string    `json:"title"`
	Description       *string   `json:"description,omitempty"`
	Venue             string    `json:"venue"`
	BannerURL         *string   `json:"banner_url,omitempty"`
	EventStartAt      time.Time `json:"event_start_at"`
	EventEndAt        time.Time `json:"event_end_at"`
	TicketSaleStartAt time.Time `json:"ticket_sale_start_at"`
	TicketSaleEndAt   time.Time `json:"ticket_sale_end_at"`
	Capacity          int       `json:"capacity"`
	Status            string    `json:"status"`
}

type EventResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       *string   `json:"description,omitempty"`
	Venue             string    `json:"venue"`
	BannerURL         *string   `json:"banner_url,omitempty"`
	EventStartAt      time.Time `json:"event_start_at"`
	EventEndAt        time.Time `json:"event_end_at"`
	TicketSaleStartAt time.Time `json:"ticket_sale_start_at"`
	TicketSaleEndAt   time.Time `json:"ticket_sale_end_at"`
	Capacity          int       `json:"capacity"`
	Status            string    `json:"status"`
}
