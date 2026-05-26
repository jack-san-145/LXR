package handlers

import (
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	//check whether the container is active or not
	active := h.Helper.IsContainerActive(con_name)
	if active {
		response.WriteJson(w, models.StartResponse{
			AlreadyActive: true,
		})
		return
	}

	//check whether the container is exists or not ,if exists make it active
	exists := h.Helper.ContainerExists(con_name)
	if exists {
		err := h.Helper.UnfreezeContainer(con_name)
		if err != nil {

			response.WriteJson(w, models.StartResponse{
				Failed: true,
			})
		} else {
			response.WriteJson(w, models.StartResponse{
				Activated: true,
			})
		}

		return
	}

	//if container doesn't exists returns response
	response.WriteJson(w, models.StartResponse{
		DoesNotExists: true,
	})

}
