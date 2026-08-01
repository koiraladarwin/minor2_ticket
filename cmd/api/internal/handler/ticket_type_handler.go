package handler

import (
	"encoding/json"
	"errors"
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

type TicketTypeHandler struct {
	service service.TicketTypeService
}

func NewTicketTypeHandler(service service.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{
		service: service,
	}
}

func (h *TicketTypeHandler) Create(w http.ResponseWriter, r *http.Request) {

	eventID, err := uuid.Parse(mux.Vars(r)["eventId"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidEventID,
		)
		return
	}

	var req dto.CreateTicketTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidRequestBody,
		)
		return
	}

	ticketType := mapper.ToTicketType(eventID, req)

	if err := h.service.Create(r.Context(), ticketType); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.Success(
		w,
		http.StatusCreated,
		response.MsgTicketTypeCreated,
		mapper.ToTicketTypeResponse(ticketType),
	)
}

func (h *TicketTypeHandler) GetByID(w http.ResponseWriter, r *http.Request) {

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
			response.MsgInvalidTicketTypeID,
		)
		return
	}

	ticketType, err := h.service.GetByIDAndUserID(
		r.Context(),
		id,
		userID,
	)

	if err != nil {

		if errors.Is(err, repository.ErrTicketTypeNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgTicketTypeNotFound,
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
		response.MsgTicketTypeFetched,
		mapper.ToTicketTypeResponse(ticketType),
	)
}

func (h *TicketTypeHandler) GetByEventID(w http.ResponseWriter, r *http.Request) {

	eventID, err := uuid.Parse(mux.Vars(r)["eventId"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidEventID,
		)
		return
	}

	ticketTypes, err := h.service.GetByEventID(
		r.Context(),
		eventID,
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
		response.MsgTicketTypesFetched,
		mapper.ToTicketTypeResponses(ticketTypes),
	)
}

func (h *TicketTypeHandler) Update(w http.ResponseWriter, r *http.Request) {

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
			response.MsgInvalidTicketTypeID,
		)
		return
	}

	var req dto.UpdateTicketTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidRequestBody,
		)
		return
	}

	ticketType, err := h.service.GetByIDAndUserID(
		r.Context(),
		id,
		userID,
	)

	if err != nil {

		if errors.Is(err, repository.ErrTicketTypeNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgTicketTypeNotFound,
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

	mapper.UpdateTicketType(ticketType, req)

	if err := h.service.Update(r.Context(), ticketType); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		response.MsgTicketTypeUpdated,
		mapper.ToTicketTypeResponse(ticketType),
	)
}

func (h *TicketTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidTicketTypeID,
		)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {

		if errors.Is(err, repository.ErrTicketTypeNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgTicketTypeNotFound,
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
		response.MsgTicketTypeDeleted,
		nil,
	)
}
