package main

import (
	"lxr-d/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/ping", h.PingHanlder)

	r.Post("/create", h.CreateHandler)
	r.Get("/run", h.RunHandler)
	r.Get("/exec", h.ExecHandler)
	r.Post("/pull_image", h.PullImageHandler)
	return r
}
