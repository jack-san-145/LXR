package helper

import (
	"log"
	"os"
	"os/exec"
)

func (h *Helper) PullImage(image string) error {
	image_env := "IMAGE=" + image
	cmd := exec.Command("../../script/pull-image.sh")
	cmd.Env = append(cmd.Environ(),
		image_env,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//run at background
	err := cmd.Start()
	if err != nil {
		log.Println("Error rootfs setup : ", err)
		return err
	}
	return nil
}
