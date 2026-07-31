package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketType struct {
	ID uuid.UUID `json:"id"`

	EventID uuid.UUID `json:"event_id"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Price float64 `json:"price"`

	Quantity  int `json:"quantity"`
	Remaining int `json:"remaining"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
