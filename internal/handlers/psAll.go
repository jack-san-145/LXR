package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

// returns all containers
func (h *Handler) PsAllHandler(w http.ResponseWriter, r *http.Request) {

	containers := h.Helper.GetAllContainers()

	response.WriteJson(w, models.PsResponse{
		Containers: containers,
	})
}
