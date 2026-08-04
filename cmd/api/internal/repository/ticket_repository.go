package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

type TicketRepository interface {

	// Purchase a ticket
	Purchase(
		ctx context.Context,
		ticketTypeID uuid.UUID,
		userID uuid.UUID,
	) (*models.Ticket, error)

	// User tickets
	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]models.TicketDetail, error)

	// Single ticket
	GetByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (*models.TicketDetail, error)
}
