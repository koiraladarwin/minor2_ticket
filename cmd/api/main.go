package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/config"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/database"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/handler"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/kafka"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/middleware"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/repository"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	kafka := kafka.NewProducer()

	authMiddleware, err := middleware.NewAuthMiddleware(cfg.PublicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Repositories
	eventRepo := repository.NewEventRepository(db)
	ticketTypeRepo := repository.NewTicketTypeRepository(db, kafka)
	ticketRepo := repository.NewTicketRepository(db)

	// Services
	eventService := service.NewEventService(eventRepo)
	ticketTypeService := service.NewTicketTypeService(ticketTypeRepo)
	ticketService := service.NewTicketService(ticketRepo, kafka)

	// Handlers
	eventHandler := handler.NewEventHandler(eventService)
	ticketTypeHandler := handler.NewTicketTypeHandler(ticketTypeService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	// Event Routes
	r.Handle(
		"/events/me",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.GetAllByUserId)),
	).Methods(http.MethodGet)

	r.Handle(
		"/events",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.GetAll)),
	).Methods(http.MethodGet)

	r.Handle(
		"/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.GetByID)),
	).Methods(http.MethodGet)

	r.Handle(
		"/eventdetails/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.GetEventDetails)),
	).Methods(http.MethodGet)

	r.Handle(
		"/events",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Create)),
	).Methods(http.MethodPost)

	r.Handle(
		"/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Update)),
	).Methods(http.MethodPut)

	r.Handle(
		"/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Delete)),
	).Methods(http.MethodDelete)

	// Ticket Type Routes

	r.Handle(
		"/events/{eventId}/ticket-types",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketTypeHandler.Create)),
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/events/{eventId}/ticket-types",
		ticketTypeHandler.GetByEventID,
	).Methods(http.MethodGet)

	r.Handle(
		"/ticket-types/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketTypeHandler.GetByID)),
	).Methods(http.MethodGet)

	r.Handle(
		"/ticket-types/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketTypeHandler.Update)),
	).Methods(http.MethodPut)

	r.Handle(
		"/ticket-types/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketTypeHandler.Delete)),
	).Methods(http.MethodDelete)

	// Ticket Routes

	r.Handle(
		"/tickets/purchase",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketHandler.Purchase)),
	).Methods(http.MethodPost)

	r.Handle(
		"/tickets/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketHandler.GetByID)),
	).Methods(http.MethodGet)

	r.Handle(
		"/me/tickets",
		authMiddleware.RequireAuth(http.HandlerFunc(ticketHandler.GetMyTickets)),
	).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Cors(r),
	}

	log.Printf("Server listening on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
