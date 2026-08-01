package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

type TicketTypeRepository interface {
	Create(ctx context.Context, ticketType *models.TicketType) error

	GetByIDAndUserID(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (*models.TicketType, error)

	GetByEventID(
		ctx context.Context,
		eventID uuid.UUID,
	) ([]models.TicketType, error)

	Update(ctx context.Context, ticketType *models.TicketType) error

	Delete(ctx context.Context, id uuid.UUID) error
}
