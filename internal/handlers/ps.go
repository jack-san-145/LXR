package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

// returns only running containers
func (h *Handler) PsHandler(w http.ResponseWriter, r *http.Request) {

	containers := h.Helper.GetActiveContainers()

	response.WriteJson(w, models.PsResponse{
		Containers: containers,
	})
}
