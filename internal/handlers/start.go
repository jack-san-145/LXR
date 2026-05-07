package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	active := h.Helper.ContainerActive(con_name)
	if active {
		response.WriteJson(w, models.StartResponse{
			AlreadyActive: true,
		})
		return
	}
	if container, ok := h.Helper.ContainerManager.AllContainers[con_name]; ok {

		err := h.Helper.ContainerSetup(container)
		if err != nil {
			response.WriteJson(w, models.StartResponse{
				Failed: true,
			})
			return
		}
		h.Helper.ContainerManager.ActiveContainers[container.ContainerName] = &container.PID

		response.WriteJson(w, models.StartResponse{
			Activated: true,
		})
		return

	}
	response.WriteJson(w, models.StartResponse{
		DoesNotExists: true,
	})
}
