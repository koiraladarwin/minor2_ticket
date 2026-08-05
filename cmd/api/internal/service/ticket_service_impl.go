package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/kafka"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
)

type ticketService struct {
	repo  repository.TicketRepository
	kafka *kafka.Producer
}

func NewTicketService(
	repo repository.TicketRepository,
	kafka *kafka.Producer,
) TicketService {
	return &ticketService{
		repo:  repo,
		kafka: kafka,
	}
}

func (s *ticketService) Purchase(
	ctx context.Context,
	ticketTypeID uuid.UUID,
	userID uuid.UUID,
) (*models.TicketDetail, error) {

	ticket, err := s.repo.Purchase(
		ctx,
		ticketTypeID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	err = s.publishTicketBought(ctx, *ticket, userID)
	if err != nil {
		fmt.Print(err.Error())
	}
	return ticket, nil
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

func (s *ticketService) publishTicketBought(
	ctx context.Context,
	ticket models.TicketDetail,
	id uuid.UUID,
) error {

	event := kafka.NewTicketCreatedEvent(
		ticket.ID.String(),
		ticket.Event.ID.String(),
		id.String(),
		ticket.Event.CreatedBy.String(),
		ticket.Status,
		ticket.QRCode,
	)

	return s.kafka.Publish(
		ctx,
		kafka.TopicTicketEvents,
		event,
	)

}

func (s *ticketService) ScanTicket(
	ctx context.Context,
	ticketID string,
	scannedBy string,
) error {
	err := s.repo.ScanTicket(ctx, ticketID, scannedBy)
	if err != nil {
		return err
	}

	event := kafka.NewTicketCheckedInEvent(
		ticketID,
		ticketID,
		scannedBy,
	)

	return s.kafka.Publish(
		ctx,
		kafka.TopicTicketEvents,
		event,
	)

}
