package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/dto"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

func ToTicketResponse(ticket *models.Ticket) dto.TicketResponse {
	return dto.TicketResponse{
		ID: ticket.ID.String(),

		EventID: ticket.EventID.String(),

		TicketTypeID: ticket.TicketTypeID.String(),

		UserID: ticket.UserID.String(),

		TicketNumber: ticket.TicketNumber,

		QRCode: ticket.QRCode,

		Status: string(ticket.Status),

		PurchasedAt: ticket.PurchasedAt,

		CheckedInAt: ticket.CheckedInAt,

		CreatedAt: ticket.CreatedAt,

		UpdatedAt: ticket.UpdatedAt,
	}
}

func ToTicketResponses(tickets []models.Ticket) []dto.TicketResponse {

	responses := make([]dto.TicketResponse, 0, len(tickets))

	for _, ticket := range tickets {
		t := ticket
		responses = append(
			responses,
			ToTicketResponse(&t),
		)
	}

	return responses
}
func ToTicketDetailResponse(ticket models.TicketDetail) dto.TicketDetailResponse {
	return dto.TicketDetailResponse{
		ID:           ticket.ID,
		TicketNumber: ticket.TicketNumber,
		QRCode:       ticket.QRCode,
		Status:       ticket.Status,
		PurchasedAt:  ticket.PurchasedAt,

		Event: struct {
			ID          uuid.UUID `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Venue       string    `json:"venue"`
			BannerURL   string    `json:"banner_url"`
			StartAt     time.Time `json:"start_at"`
			EndAt       time.Time `json:"end_at"`
			Status      string    `json:"status"`
		}{
			ID:          ticket.Event.ID,
			Title:       ticket.Event.Title,
			Description: ticket.Event.Description,
			Venue:       ticket.Event.Venue,
			BannerURL:   ticket.Event.BannerURL,
			StartAt:     ticket.Event.StartAt,
			EndAt:       ticket.Event.EndAt,
			Status:      ticket.Event.Status,
		},

		TicketType: struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Price float64   `json:"price"`
		}{
			ID:    ticket.TicketType.ID,
			Name:  ticket.TicketType.Name,
			Price: ticket.TicketType.Price,
		},
	}
}
func ToTicketDetailResponseList(
	tickets []models.TicketDetail,
) []dto.TicketDetailResponse {

	responses := make([]dto.TicketDetailResponse, 0, len(tickets))

	for _, ticket := range tickets {
		responses = append(
			responses,
			ToTicketDetailResponse(ticket),
		)
	}

	return responses
}
