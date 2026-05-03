package helper

import (
	"fmt"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func (h *Helper) RootfsSetup(con *models.Container) error {
	id := uuid.New()
	con.ContainerId = strings.Join(strings.Split(id.String(), "-"), "")

	container_name_env := "CONTAINER_NAME=" + con.ContainerName
	image_name_env := "IMAGE_NAME=" + con.Image
	container_id_env := "CONTAINER_ID=" + con.ContainerId

	//to check the image locally
	exists := h.CheckImageLocally(con.Image)
	if !exists {
		err := h.PullImage(con.Image)
		if err != nil {
			return err
		}
	}

	//run the script with env
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
