package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"lxr-d/internal/models"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func (h *LXRHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var con models.NewContainer
	err := json.NewDecoder(r.Body).Decode(&con)
	if err != nil {
		log.Println("Creation Error: ", err)
		return
	}

	RootfsSetup(con)
}

func RootfsSetup(con models.NewContainer) {
	id := uuid.New()
	container_id := strings.Join(strings.Split(id.String(), "-"), "")

	container_name_env := "CONTAINER_NAME=" + con.ContainerName
	image_name_env := "IMAGE_NAME=" + con.Image
	container_id_env := "CONTAINER_ID=" + container_id

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
	}

}
