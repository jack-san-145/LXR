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
	img_name_env := "IMAGE_NAME=" + container.Image
	pass_env := "PASSWORD=" + container.ContainerName

	cmd := exec.Command(
		"unshare",
		"--pid",
		"--mount",
		"--uts",
		"--net",
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
		img_name_env,
		pass_env)

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
