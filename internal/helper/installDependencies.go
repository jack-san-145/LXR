package helper

import (
	"log"
	"lxr-d/internal/models"
	"os"
	"os/exec"
	"strconv"
)

func (h *Helper) InstallDependencies(con *models.Container) {

	pass_env := "PASSWORD=" + con.ContainerName

	cmd := exec.Command(
		"nsenter",
		"--target", strconv.Itoa(con.PID),
		"--pid", "--mount", "--uts", "--net",
		"bash", "../../script/install-dependencies.sh",
	)

}
