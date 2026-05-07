package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StopHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	exists, err := h.Helper.StopProcess(con_name)
	if !exists || err != nil {
		response.WriteJson(w, models.StopResponse{
			Exists:  exists,
			Stopped: false,
		})
	}

	response.WriteJson(w, models.StopResponse{
		Exists:  exists,
		Stopped: true,
	})
}
