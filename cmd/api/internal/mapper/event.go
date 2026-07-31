package mapper

import (
	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/dto"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

func ToEvent(req dto.CreateEventRequest, userID uuid.UUID) *models.Event {
	return &models.Event{
		Title:             req.Title,
		Description:       req.Description,
		Venue:             req.Venue,
		BannerURL:         req.BannerURL,
		EventStartAt:      req.EventStartAt,
		EventEndAt:        req.EventEndAt,
		TicketSaleStartAt: req.TicketSaleStartAt,
		TicketSaleEndAt:   req.TicketSaleEndAt,
		Capacity:          req.Capacity,
		CreatedBy:         userID,
	}
}

func UpdateEvent(event *models.Event, req dto.UpdateEventRequest) {
	event.Title = req.Title
	event.Description = req.Description
	event.Venue = req.Venue
	event.BannerURL = req.BannerURL
	event.EventStartAt = req.EventStartAt
	event.EventEndAt = req.EventEndAt
	event.TicketSaleStartAt = req.TicketSaleStartAt
	event.TicketSaleEndAt = req.TicketSaleEndAt
	event.Capacity = req.Capacity
}

func ToEventResponse(event *models.Event) dto.EventResponse {
	return dto.EventResponse{
		ID:                event.ID.String(),
		Title:             event.Title,
		Description:       event.Description,
		Venue:             event.Venue,
		BannerURL:         event.BannerURL,
		EventStartAt:      event.EventStartAt,
		EventEndAt:        event.EventEndAt,
		TicketSaleStartAt: event.TicketSaleStartAt,
		TicketSaleEndAt:   event.TicketSaleEndAt,
		Capacity:          event.Capacity,
		Status:            string(event.Status),
	}
}

func ToEventResponses(events []models.Event) []dto.EventResponse {
	response := make([]dto.EventResponse, 0, len(events))

	for _, event := range events {
		response = append(response, ToEventResponse(&event))
	}

	return response
}
