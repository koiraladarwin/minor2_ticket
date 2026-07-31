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
