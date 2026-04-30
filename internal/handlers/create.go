package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var con models.Container
	err := json.NewDecoder(r.Body).Decode(&con)
	if err != nil {
		log.Println("Creation Error: ", err)
		return
	}
	err = RootfsSetup(&con)
	if err != nil {
		log.Println("Error during RootfsSetup: ", err)
		response.WriteJson(w, models.CreationResponse{IsCreated: false})
		return
	}
	response.WriteJson(w, models.CreationResponse{
		IsCreated:     true,
		ContainerName: con.ContainerName,
		ContainerId:   con.ContainerId,
	})

	h.Helper.ContainerManager.AllContainers[con.ContainerName] = &con

}

func RootfsSetup(con *models.Container) error {
	id := uuid.New()
	con.ContainerId = strings.Join(strings.Split(id.String(), "-"), "")

	container_name_env := "CONTAINER_NAME=" + con.ContainerName
	image_name_env := "IMAGE_NAME=" + con.Image
	container_id_env := "CONTAINER_ID=" + con.ContainerId

	cmd := exec.Command("../../script/rootfs-setup.sh")
	cmd.Env = append(os.Environ(),
		container_name_env,
		container_id_env,
		image_name_env)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error rootfs setup : ", err)
		return err
	}
	return nil
}
