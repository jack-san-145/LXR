package helper

import (
	"io"
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	// "strconv"
)

func (h *Helper) InstallDependencies(cb *models.ContainerBuilder) error {

	//set image image as password for container initially

	img_name_env := "IMAGE_NAME=" + cb.Container.Image

	cmd := exec.Command(
		"sudo",
		"-E",
		"unshare",
		"--pid",
		"--mount",
		"--fork",
		"--",
		"bash",
		"../../script/install-dependencies.sh",
	)

	//inject env to the script
	cmd.Env = append(os.Environ(),
		img_name_env,
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
		log.Println("Error container dependencies: ", err)
		return err
	}

	//wait until script to complete
	err = cmd.Wait()
	if err != nil {
		log.Println("Install dependency error: ", err)
		return err
	}

	return nil

}
