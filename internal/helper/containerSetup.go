package helper

import (
	"fmt"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"syscall"
)

func (h *Helper) ContainerSetup(container *models.Container) error {

	container_name_env := "CONTAINER_NAME=" + container.ContainerName
	container_id_env := "CONTAINER_ID=" + container.ContainerId
	img_name_env := "IMAGE=" + container.Image

	cmd := exec.Command(
		"unshare",
		"--pid",
		"--mount",
		"--uts",
		"--map-root-user",
		"--fork",
		"--",
		"bash",
		"../../script/container-setup.sh",
	)

	//switch non-root to spawn containers
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 1000, // non-root
			Gid: 1000,
		},
	}

	//inject env to the script
	cmd.Env = append(os.Environ(),
		container_name_env,
		container_id_env,
		img_name_env)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//run script in background
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error container setup : ", err)
		return err
	}
	container.PID = cmd.Process.Pid
	fmt.Println("container Pid = ", container.PID)

	// waits for child exit signal and reaps the process
	go func() {
		cmd.Wait()
	}()
	return nil
}
