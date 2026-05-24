package helper

import (
	"github.com/google/uuid"
	"io"
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strings"
)

func (h *Helper) RootfsSetup(conCreater *models.ContainerCreater) error {

	id := uuid.New()
	conCreater.Container.ContainerId = strings.Join(strings.Split(id.String(), "-"), "")

	container_name_env := "CONTAINER_NAME=" + conCreater.Container.ContainerName
	image_name_env := "IMAGE_NAME=" + conCreater.Container.Image
	container_id_env := "CONTAINER_ID=" + conCreater.Container.ContainerId

	//run the script with env
	cmd := exec.Command("../../script/rootfs-setup.sh")
	cmd.Env = append(os.Environ(),
		container_name_env,
		container_id_env,
		image_name_env)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Get a pipe to read the command's standard output and error (stdout and stderr)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Stream stdoud(normal output) from the command to the client
	go func() {
		io.Copy(conCreater.Conn, stdout)
	}()

	// Stream stderr(error/warning output) from the command to the client
	go func() {
		io.Copy(conCreater.Conn, stderr)
	}()

	//start to run script at background
	err := cmd.Start()
	if err != nil {
		log.Println("rootfs setup Error: ", err)
		return err
	}

	//wait until script to complete
	err = cmd.Wait()
	if err != nil {
		log.Println("rootfs setup Error: ", err)
		return err
	}

	return nil

}
