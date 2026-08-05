package kafka

import "time"

type TicketCache struct {
	EventType    string    `json:"event_type"`
	TicketID     string    `json:"ticket_id"`
	QrCode       string    `json:"qrcode"`
	EventID      string    `json:"event_id"`
	BuyerID      string    `json:"buyer_id"`
	EventOwnerID string    `json:"event_owner_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
