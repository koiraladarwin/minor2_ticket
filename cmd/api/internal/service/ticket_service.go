package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

type TicketService interface {

	// Purchase a ticket
	Purchase(
		ctx context.Context,
		ticketTypeID uuid.UUID,
		userID uuid.UUID,
	) (*models.TicketDetail, error)

	// Get one ticket owned by user
	GetByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (*models.TicketDetail, error)

	// Get all tickets owned by user
	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]models.TicketDetail, error)

	//Scan ticket of a user
	ScanTicket(
		ctx context.Context,
		ticketID string,
		scannedBy string,
	) error
}
