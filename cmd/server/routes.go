package main

import (
	"lxr-d/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/ping", h.PingHanlder)

	r.Post("/create", h.CreateHandler)
	r.Post("/start", h.StartHandler)
	r.Get("/stop", h.StopHandler)
	r.Get("/exec", h.ExecHandler)
	r.Delete("/kill", h.KillHanlder)
	r.Post("/pull_image", h.PullImageHandler)

	r.Get("/ps", h.PsHandler)
	r.Get("/ps/all", h.PsAllHandler)

	return r
}
