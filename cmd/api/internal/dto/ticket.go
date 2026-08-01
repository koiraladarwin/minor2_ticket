package dto

import "time"

type PurchaseTicketRequest struct {
	TicketTypeID string `json:"ticket_type_id"`
}

type TicketResponse struct {
	ID string `json:"id"`

	EventID string `json:"event_id"`

	TicketTypeID string `json:"ticket_type_id"`

	UserID string `json:"user_id"`

	TicketNumber string `json:"ticket_number"`

	QRCode string `json:"qr_code"`

	Status string `json:"status"`

	PurchasedAt time.Time `json:"purchased_at"`

	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
