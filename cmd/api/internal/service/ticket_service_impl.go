package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
)

type ticketService struct {
	repo repository.TicketRepository
}

func NewTicketService(
	repo repository.TicketRepository,
) TicketService {
	return &ticketService{
		repo: repo,
	}
}

func (s *ticketService) Purchase(
	ctx context.Context,
	ticketTypeID uuid.UUID,
	userID uuid.UUID,
) (*models.Ticket, error) {
	return s.repo.Purchase(
		ctx,
		ticketTypeID,
		userID,
	)
}

func (s *ticketService) GetByIDAndUserID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (*models.TicketDetail, error) {
	return s.repo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)
}

func (s *ticketService) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.TicketDetail, error) {
	return s.repo.GetByUserID(
		ctx,
		userID,
	)
}
