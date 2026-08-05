package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/dto"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/mapper"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/middleware"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/response"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/service"
)

type TicketHandler struct {
	service service.TicketService
}

func NewTicketHandler(service service.TicketService) *TicketHandler {
	return &TicketHandler{
		service: service,
	}
}

func (h *TicketHandler) Purchase(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.PurchaseTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidRequestBody,
		)
		return
	}

	ticketTypeID, err := uuid.Parse(req.TicketTypeID)
	if err != nil {
		log.Print(err)
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidTicketID,
		)
		return
	}

	userID, err := uuid.Parse(middleware.UserID(r.Context()))
	if err != nil {
		response.Error(
			w,
			http.StatusUnauthorized,
			response.MsgUnauthorized,
		)
		return
	}

	ticket, err := h.service.Purchase(
		r.Context(),
		ticketTypeID,
		userID,
	)
	if err != nil {

		switch {

		case errors.Is(err, repository.ErrAlreadyPurchased):
			response.Error(
				w,
				http.StatusConflict,
				err.Error(),
			)

		case errors.Is(err, repository.ErrTicketSoldOut):
			response.Error(
				w,
				http.StatusConflict,
				err.Error(),
			)

		case errors.Is(err, repository.ErrTicketSalesClosed):
			response.Error(
				w,
				http.StatusBadRequest,
				err.Error(),
			)

		case errors.Is(err, repository.ErrEventUnavailable):
			response.Error(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		default:
			response.Error(
				w,
				http.StatusInternalServerError,
				response.MsgInternalServer,
			)
		}

		return
	}

	response.Success(
		w,
		http.StatusCreated,
		"Ticket purchased successfully",
		mapper.ToTicketDetailResponse(*ticket),
	)
}

func (h *TicketHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, err := uuid.Parse(middleware.UserID(r.Context()))
	if err != nil {
		response.Error(
			w,
			http.StatusUnauthorized,
			response.MsgUnauthorized,
		)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidTicketID,
		)
		return
	}

	ticket, err := h.service.GetByIDAndUserID(
		r.Context(),
		id,
		userID,
	)
	if err != nil {

		if errors.Is(err, repository.ErrTicketNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgTicketNotFound,
			)
			return
		}

		response.Error(
			w,
			http.StatusInternalServerError,
			response.MsgInternalServer,
		)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		response.MsgTicketFetched,
		mapper.ToTicketDetailResponse(*ticket),
	)
}

func (h *TicketHandler) GetMyTickets(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, err := uuid.Parse(middleware.UserID(r.Context()))
	if err != nil {
		response.Error(
			w,
			http.StatusUnauthorized,
			response.MsgUnauthorized,
		)
		return
	}

	tickets, err := h.service.GetByUserID(
		r.Context(),
		userID,
	)
	if err != nil {
		response.Error(
			w,
			http.StatusInternalServerError,
			response.MsgInternalServer,
		)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		response.MsgTicketsFetched,
		mapper.ToTicketDetailResponseList(tickets),
	)
}
