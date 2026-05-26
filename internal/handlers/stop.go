package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StopHandler(w http.ResponseWriter, r *http.Request) {

	containerName := r.URL.Query().Get("container_name")

	//check whether the container is active or not
	active := h.Helper.IsContainerActive(containerName)
	if active {
		h.Helper.FreezeContainer(containerName)
		response.WriteJson(w, models.StopResponse{
			Exists:  true,
			Stopped: true,
		})
		return
	}

	//check whether the container is exists or not
	exists := h.Helper.ContainerExists(containerName)
	if exists {
		response.WriteJson(w, models.StopResponse{
			Exists:  true,
			Stopped: false,
		})
		return
	}

	//if container doesn't exists returns response
	response.WriteJson(w, models.StopResponse{
		Exists:  false,
		Stopped: false,
	})

}
