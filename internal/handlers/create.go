package handlers

import (
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var con models.Container
	err := json.NewDecoder(r.Body).Decode(&con)
	if err != nil {
		log.Println("Creation Error: ", err)
		return
	}

	//check the container already exists or not
	exists := h.Helper.ContainerExists(con.ContainerName)
	if exists {
		response.WriteJson(w, models.CreationResponse{
			IsCreated:     false,
			AlreadyExists: true,
		})
		return
	}

	//if not already exists creats new container environment (only setup containers rootfs)
	err = h.Helper.RootfsSetup(&con)
	if err != nil {
		log.Println("Error during RootfsSetup: ", err)
		response.WriteJson(w, models.CreationResponse{IsCreated: false})
		return
	}

	//send container creattion response to client
	response.WriteJson(w, models.CreationResponse{
		IsCreated:     true,
		ContainerName: con.ContainerName,
		ContainerId:   con.ContainerId,
	})

	//add newly created container to allContainers
	h.Helper.ContainerManager.AllContainers[con.ContainerName] = &con

}
