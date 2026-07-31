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

type EventHandler struct {
	service service.EventService
}

func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidRequestBody,
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

	event := mapper.ToEvent(req, userID)

	if err := h.service.Create(r.Context(), event); err != nil {
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
		response.MsgEventCreated,
		mapper.ToEventResponse(event),
	)
}

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidEventID,
		)
		return
	}

	event, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgEventNotFound,
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
		response.MsgEventFetched,
		mapper.ToEventResponse(event),
	)
}

func (h *EventHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Println(err)
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
		response.MsgEventsFetched,
		mapper.ToEventResponses(events),
	)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidEventID,
		)
		return
	}

	var req dto.UpdateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidRequestBody,
		)
		return
	}

	event, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgEventNotFound,
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

	mapper.UpdateEvent(event, req)

	if err := h.service.Update(r.Context(), event); err != nil {
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
		response.MsgEventUpdated,
		mapper.ToEventResponse(event),
	)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			response.MsgInvalidEventID,
		)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			response.Error(
				w,
				http.StatusNotFound,
				response.MsgEventNotFound,
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
		response.MsgEventDeleted,
		nil,
	)
}
