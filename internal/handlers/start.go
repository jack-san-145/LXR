package handlers

import (
	"log"
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	//check whether the container is active or not
	active := h.Helper.ContainerActive(con_name)
	if active {
		response.WriteJson(w, models.StartResponse{
			AlreadyActive: true,
		})
		return
	}

	//check whether the container is exists or not
	if con, ok := h.Helper.ContainerManager.AllContainers[con_name]; ok {

		//start container with rootfs
		err := h.Helper.ContainerSetup(con)
		if err != nil {
			response.WriteJson(w, models.StartResponse{
				Failed: true,
			})
			return
		}
		h.Helper.ContainerManager.ActiveContainers[con.ContainerName] = &con.PID

		//setup networking for container
		err = h.Helper.SetupContainerNetworking(con)
		if err != nil {
			log.Println("Error during Container Networking: ", err)
			response.WriteJson(w, models.StartResponse{
				Failed: true,
			})
			return
		}

		//returns activated response
		response.WriteJson(w, models.StartResponse{
			Activated: true,
		})
		return

	}

	//if container doesn't exists returns response
	response.WriteJson(w, models.StartResponse{
		DoesNotExists: true,
	})
}
