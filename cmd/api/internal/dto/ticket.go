package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type TicketDetailResponse struct {
	ID           uuid.UUID `json:"id"`
	TicketNumber string    `json:"ticket_number"`
	QRCode       string    `json:"qr_code"`
	Status       string    `json:"status"`

	PurchasedAt time.Time `json:"purchased_at"`

	Event struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Venue       string    `json:"venue"`
		BannerURL   string    `json:"banner_url"`

		StartAt time.Time `json:"start_at"`
		EndAt   time.Time `json:"end_at"`

		Status string `json:"status"`
	} `json:"event"`

	TicketType struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Price float64   `json:"price"`
	} `json:"ticket_type"`
}

type ScanTicketRequest struct {
	TicketID string `json:"ticket_id"`
}
