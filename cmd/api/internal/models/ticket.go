package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketActive    TicketStatus = "ACTIVE"
	TicketUsed      TicketStatus = "USED"
	TicketCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID uuid.UUID `json:"id"`

	EventID uuid.UUID `json:"event_id"`

	TicketTypeID uuid.UUID `json:"ticket_type_id"`

	UserID uuid.UUID `json:"user_id"`

	TicketNumber string `json:"ticket_number"`

	QRCode string `json:"qr_code"`

	Status TicketStatus `json:"status"`

	PurchasedAt time.Time  `json:"purchased_at"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TicketDetail struct {
	ID           uuid.UUID `json:"id"`
	TicketNumber string    `json:"ticket_number"`
	QRCode       string    `json:"qr_code"`
	Status       string    `json:"status"`

	PurchasedAt time.Time  `json:"purchased_at"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`

	Event struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Venue       string    `json:"venue"`
		BannerURL   string    `json:"banner_url"`

		StartAt time.Time `json:"start_at"`
		EndAt   time.Time `json:"end_at"`

		Status    string    `json:"status"`
		CreatedBy uuid.UUID `json:"created_by"`
	} `json:"event"`

	TicketType struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
	} `json:"ticket_type"`

	UserID uuid.UUID `json:"user_id"`
}
