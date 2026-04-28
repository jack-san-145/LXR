package handlers

import (
	"fmt"
	"lxr-d/internal/models"
	"lxr-d/internal/response"
	"net/http"
	"os"
	"os/exec"
)

func (h *LXRHandler) RunHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	if container, ok := h.Helper.ContainerManager.AllContainers[con_name]; ok {
		err := ContainerSetup(container)
		if err != nil {
			response.WriteJson(w, map[string]bool{"active": false})
			return
		}
		h.Helper.ContainerManager.ActiveContainers[container.ContainerName] = &container.PID
		response.WriteJson(w, map[string]bool{"active": true})
	}
}

func ContainerSetup(container *models.Container) error {

	container_name_env := "CONTAINER_NAME=" + container.ContainerName
	container_id_env := "CONTAINER_ID=" + container.ContainerId

	cmd := exec.Command(
		"unshare",
		"--pid",
		"--mount",
		"--map-root-user",
		"--fork",
		"--",
		"bash",
		"../../script/container-setup.sh",
	)

	cmd.Env = append(os.Environ(),
		container_name_env,
		container_id_env)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error container setup : ", err)
		return err
	}
	container.PID = cmd.Process.Pid
	fmt.Println("container Pid = ", container.PID)

	return nil
}
