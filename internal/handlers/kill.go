package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) KillHanlder(w http.ResponseWriter, r *http.Request) {

	conName := r.URL.Query().Get("container_name")
	if con, ok := h.Helper.ContainerManager.AllContainers[conName]; ok {
		err := h.Helper.KillContainer(con)
		if err != nil {
			response.WriteJson(w, models.KillResponse{
				Exists: true,
				Killed: false,
			})
			return
		}

		response.WriteJson(w, models.KillResponse{
			Exists: true,
			Killed: true,
		})
		return
	}
	response.WriteJson(w, models.KillResponse{
		Exists: false,
		Killed: false,
	})

}
