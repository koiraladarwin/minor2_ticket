package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
)

var (
	ErrInvalidEventDates    = errors.New("event end must be after event start")
	ErrInvalidSaleDates     = errors.New("invalid ticket sale dates")
	ErrSaleAfterEventStarts = errors.New("ticket sales must end before the event starts")
	ErrInvalidCapacity      = errors.New("capacity must be greater than zero")
)

type EventService interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetAll(ctx context.Context) ([]models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{
		repo: repo,
	}
}

func (s *eventService) Create(ctx context.Context, event *models.Event) error {

	if event.Capacity <= 0 {
		return ErrInvalidCapacity
	}

	if !event.EventEndAt.After(event.EventStartAt) {
		return ErrInvalidEventDates
	}

	if !event.TicketSaleEndAt.After(event.TicketSaleStartAt) {
		return ErrInvalidSaleDates
	}

	if event.TicketSaleEndAt.After(event.EventStartAt) {
		return ErrSaleAfterEventStarts
	}

	if event.Status == "" {
		event.Status = models.EventDraft
	}

	return s.repo.Create(ctx, event)
}

func (s *eventService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Event, error) {

	return s.repo.GetByID(ctx, id)
}

func (s *eventService) GetAll(
	ctx context.Context,
) ([]models.Event, error) {

	return s.repo.GetAll(ctx)
}

func (s *eventService) Update(
	ctx context.Context,
	event *models.Event,
) error {

	if event.Capacity <= 0 {
		return ErrInvalidCapacity
	}

	if !event.EventEndAt.After(event.EventStartAt) {
		return ErrInvalidEventDates
	}

	if !event.TicketSaleEndAt.After(event.TicketSaleStartAt) {
		return ErrInvalidSaleDates
	}

	if event.TicketSaleEndAt.After(event.EventStartAt) {
		return ErrSaleAfterEventStarts
	}

	return s.repo.Update(ctx, event)
}

func (s *eventService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	return s.repo.Delete(ctx, id)
}
