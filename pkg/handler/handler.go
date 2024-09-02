package handler

import (
	"medodstest/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/sign-up", h.signUp)
		mux.Post("/sign-in", h.signIn)
		mux.With(h.userIdentity).Post("/refresh", h.refresh)
	})

	return mux
}
