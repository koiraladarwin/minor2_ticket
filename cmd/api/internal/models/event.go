package models

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventDraft     EventStatus = "DRAFT"
	EventPublished EventStatus = "PUBLISHED"
	EventCancelled EventStatus = "CANCELLED"
	EventCompleted EventStatus = "COMPLETED"
)

type Event struct {
	ID uuid.UUID `json:"id"`

	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Venue       string  `json:"venue"`
	BannerURL   *string `json:"banner_url,omitempty"`

	EventStartAt      time.Time `json:"event_start_at"`
	EventEndAt        time.Time `json:"event_end_at"`
	TicketSaleStartAt time.Time `json:"ticket_sale_start_at"`
	TicketSaleEndAt   time.Time `json:"ticket_sale_end_at"`

	Capacity int         `json:"capacity"`
	Status   EventStatus `json:"status"`

	CreatedBy uuid.UUID `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
