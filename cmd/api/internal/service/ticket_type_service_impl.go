package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
)

type ticketTypeService struct {
	repo repository.TicketTypeRepository
}

func NewTicketTypeService(
	repo repository.TicketTypeRepository,
) TicketTypeService {
	return &ticketTypeService{
		repo: repo,
	}
}

func (s *ticketTypeService) Create(
	ctx context.Context,
	ticketType *models.TicketType,
) error {

	if ticketType.Name == "" {
		return errors.New("ticket type name is required")
	}

	if ticketType.Price < 0 {
		return errors.New("ticket price cannot be negative")
	}

	if ticketType.Quantity <= 0 {
		return errors.New("ticket quantity must be greater than zero")
	}

	ticketType.Remaining = ticketType.Quantity

	return s.repo.Create(ctx, ticketType)
}

func (s *ticketTypeService) GetByIDAndUserID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (*models.TicketType, error) {

	return s.repo.GetByIDAndUserID(
		ctx,
		id,
		userID,
	)
}

func (s *ticketTypeService) GetByEventID(
	ctx context.Context,
	eventID uuid.UUID,
) ([]models.TicketType, error) {

	return s.repo.GetByEventID(
		ctx,
		eventID,
	)
}

func (s *ticketTypeService) Update(
	ctx context.Context,
	ticketType *models.TicketType,
) error {

	if ticketType.Name == "" {
		return errors.New("ticket type name is required")
	}

	if ticketType.Price < 0 {
		return errors.New("ticket price cannot be negative")
	}

	if ticketType.Quantity <= 0 {
		return errors.New("ticket quantity must be greater than zero")
	}

	if ticketType.Remaining > ticketType.Quantity {
		return errors.New("remaining tickets cannot exceed quantity")
	}

	return s.repo.Update(
		ctx,
		ticketType,
	)
}

func (s *ticketTypeService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	return s.repo.Delete(
		ctx,
		id,
	)
}
