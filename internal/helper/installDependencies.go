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

	//inject env to the script
	cmd.Env = append(os.Environ(),
		pass_env,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//run script in background
	err := cmd.Start()
	if err != nil {
		log.Println("Error container dependencies: ", err)
	}

}
