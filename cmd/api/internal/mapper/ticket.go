package mapper

import (
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
