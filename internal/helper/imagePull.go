package helper

import (
	"io"
	"log"
	"lxr-d/internal/models"
	"os/exec"
)

func (h *Helper) PullImage(cb *models.ContainerBuilder) error {

	image_env := "IMAGE_NAME=" + cb.Container.Image

	cmd := exec.Command("../../script/pull-image.sh")
	cmd.Env = append(cmd.Environ(),
		image_env,
	)

	// Get a pipe to read the command's standard output and error (stdout and stderr)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Stream stdoud(normal output) from the command to the client
	go func() {
		io.Copy(cb.Conn, stdout)
	}()

	// Stream stderr(error/warning output) from the command to the client
	go func() {
		io.Copy(cb.Conn, stderr)
	}()

	//start to run script at background
	err := cmd.Start()
	if err != nil {
		log.Println("Image script Error : ", err)
		return err
	}

	//wait until script to complete
	err = cmd.Wait()
	if err != nil {
		log.Println("Image pull Error: ", err)
		return err
	}
	return nil
}
