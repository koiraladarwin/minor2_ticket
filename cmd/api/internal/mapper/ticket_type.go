package mapper

import (
	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/dto"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

func ToTicketType(eventID uuid.UUID, req dto.CreateTicketTypeRequest) *models.TicketType {
	return &models.TicketType{
		EventID:     eventID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Remaining:   req.Quantity,
	}
}

func UpdateTicketType(tt *models.TicketType, req dto.UpdateTicketTypeRequest) {
	tt.Name = req.Name
	tt.Description = req.Description
	tt.Price = req.Price
	tt.Quantity = req.Quantity
}

func ToTicketTypeResponse(tt *models.TicketType) dto.TicketTypeResponse {
	return dto.TicketTypeResponse{
		ID:          tt.ID,
		EventID:     tt.EventID,
		Name:        tt.Name,
		Description: tt.Description,
		Price:       tt.Price,
		Quantity:    tt.Quantity,
		Remaining:   tt.Remaining,
	}
}

func ToTicketTypeResponses(types []models.TicketType) []dto.TicketTypeResponse {
	res := make([]dto.TicketTypeResponse, 0, len(types))

	for _, t := range types {
		res = append(res, ToTicketTypeResponse(&t))
	}

	return res
}
