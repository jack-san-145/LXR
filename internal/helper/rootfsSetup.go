package helper

import (
	"fmt"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func (h *Helper) RootfsSetup(conCreater *models.ContainerCreater) error {

	id := uuid.New()
	conCreater.Container.ContainerId = strings.Join(strings.Split(id.String(), "-"), "")

	container_name_env := "CONTAINER_NAME=" + conCreater.Container.ContainerName
	image_name_env := "IMAGE_NAME=" + conCreater.Container.Image
	container_id_env := "CONTAINER_ID=" + conCreater.Container.ContainerId

	// //to check the image locally
	// exists := h.CheckImageLocally(conCreater.Container.Image)
	// if !exists {
	// 	_, err := h.PullImage(conCreater.Container.Image)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
