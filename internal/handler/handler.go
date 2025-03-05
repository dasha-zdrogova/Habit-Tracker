package handler

import (
	"habit-tracker/internal/handler/auth"
	"habit-tracker/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	services     *service.Services
	tokenManager *auth.TokenManager
}

func NewHandler(sevices *service.Services) *Handler {
	return &Handler{
		services:     sevices,
		tokenManager: auth.NewTokenManager(),
	}
}

func (h *Handler) Register(router chi.Router) {
	//TODO: добавить логи
	router.Use(middleware.Recoverer)

	// Аутентификация
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.register)
		r.Post("/login", h.login)
	})

	// API с аутентификацией
	router.Route("/api", func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Route("/habit", func(r chi.Router) {
			r.Post("/create", h.createHabit)

			r.Get("/info/{id}", h.WithHabitAccess(h.getHabitInfo))
			r.Put("/mark/{id}", h.WithHabitAccess(h.markHabit))
			r.Delete("/delete/{id}", h.WithHabitAccess(h.deleteHabit))

			// r.Get("/all", h.getHabits)
		})
	})

}
