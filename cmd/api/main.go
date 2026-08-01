package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/config"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/database"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/handler"
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

	authMiddleware, err := middleware.NewAuthMiddleware(cfg.PublicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)
	r.HandleFunc("/events", eventHandler.GetAll).Methods(http.MethodGet)

	r.Handle("/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.GetByID)),
	).Methods(http.MethodGet)

	r.Handle("/events",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Create)),
	).Methods(http.MethodPost)
	r.Handle("/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Update)),
	).Methods(http.MethodPut)
	r.Handle("/events/{id}",
		authMiddleware.RequireAuth(http.HandlerFunc(eventHandler.Delete)),
	).Methods(http.MethodDelete)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("Server listening on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
