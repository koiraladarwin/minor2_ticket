package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/dto"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/kafka"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/utils"
)

var (
	ErrInvalidQR              = errors.New("invalid qr")
	ErrQRExpired              = errors.New("qr expired")
	ErrInvalidTicketSignature = errors.New("invalid ticket signature")
)

type ticketService struct {
	repo  repository.TicketRepository
	redis *redis.Client
	kafka *kafka.Producer
}

func NewTicketService(
	repo repository.TicketRepository,
	kafka *kafka.Producer,
	redis *redis.Client,
) TicketService {
	return &ticketService{
		repo:  repo,
		kafka: kafka,
		redis: redis,
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
	expTime int64,
	signature string,
	scannedBy string,
) error {

	if time.Now().Unix() > expTime {

		return ErrQRExpired

	}

	cacheKey := "ticket:" + ticketID

	data, err := s.redis.Get(
		ctx,
		cacheKey,
	).Result()

	if err != nil {

		return repository.ErrTicketNotFound

	}

	var ticketCache dto.TicketCache

	err = json.Unmarshal(
		[]byte(data),
		&ticketCache,
	)

	if err != nil {

		return err

	}

	valid := utils.VerifyHMAC(
		ticketID,
		expTime,
		signature,
		ticketCache.QrCode,
	)

	if !valid {

		return ErrInvalidTicketSignature

	}

	ticket, err := s.repo.GetByID(
		ctx,
		uuid.MustParse(ticketID),
	)

	if err != nil {

		return err

	}

	err = s.repo.ScanTicket(
		ctx,
		ticketID,
		scannedBy,
	)

	if err != nil {

		return err

	}

	event := kafka.NewTicketCheckedInEvent(
		ticketID,
		ticket.Event.ID.String(),
		scannedBy,
	)

	return s.kafka.Publish(
		ctx,
		kafka.TopicTicketEvents,
		event,
	)

}
