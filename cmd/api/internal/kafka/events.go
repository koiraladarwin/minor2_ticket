package kafka

import "time"

const (
	TopicTicketEvents = "ticket.events"

	EventTicketCreated   = "ticket.created"
	EventTicketCheckedIn = "ticket.checked_in"
	EventTicketCancelled = "ticket.cancelled"
)

type TicketCreatedEvent struct {
	EventType    string    `json:"event_type"`
	TicketID     string    `json:"ticket_id"`
	EventID      string    `json:"event_id"`
	BuyerID      string    `json:"buyer_id"`
	EventOwnerID string    `json:"event_owner_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type TicketCheckedInEvent struct {
	EventType string    `json:"event_type"`
	TicketID  string    `json:"ticket_id"`
	EventID   string    `json:"event_id"`
	CheckedBy string    `json:"checked_by"`
	CheckedAt time.Time `json:"checked_at"`
}

type TicketCancelledEvent struct {
	EventType   string    `json:"event_type"`
	TicketID    string    `json:"ticket_id"`
	EventID     string    `json:"event_id"`
	Reason      string    `json:"reason"`
	CancelledAt time.Time `json:"cancelled_at"`
}

func NewTicketCreatedEvent(
	ticketID string,
	eventID string,
	buyerID string,
	eventOwnerID string,
	status string,
) TicketCreatedEvent {

	return TicketCreatedEvent{
		EventType:    EventTicketCreated,
		TicketID:     ticketID,
		EventID:      eventID,
		BuyerID:      buyerID,
		EventOwnerID: eventOwnerID,
		Status:       status,
		CreatedAt:    time.Now(),
	}
}

func NewTicketCheckedInEvent(
	ticketID string,
	eventID string,
	checkedBy string,
) TicketCheckedInEvent {

	return TicketCheckedInEvent{
		EventType: EventTicketCheckedIn,
		TicketID:  ticketID,
		EventID:   eventID,
		CheckedBy: checkedBy,
		CheckedAt: time.Now(),
	}
}

func NewTicketCancelledEvent(
	ticketID string,
	eventID string,
	reason string,
) TicketCancelledEvent {

	return TicketCancelledEvent{
		EventType:   EventTicketCancelled,
		TicketID:    ticketID,
		EventID:     eventID,
		Reason:      reason,
		CancelledAt: time.Now(),
	}
}
