package main

import (
	"lxr-d/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(lxr *handlers.LXRHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/ping", lxr.PingHanlder)

	r.Post("/create", lxr.CreateHandler)
	r.Get("/run", lxr.RunHandler)
	r.Get("/exec", lxr.ExecHandler)
	return r
}
