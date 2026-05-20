package helper

import (
	"fmt"
	"log"
	"os/exec"
)

// to get container init PID
func (h *Helper) GetContainerInitPid(conPID int) (string, error) {

	//extract container's init pid with parent pid
	CMD := fmt.Sprintf("ps -ef | grep %v | head -n 2 | tail -n 1 | awx '{print $2}' ", conPID)

	cmd := exec.Command("bash", "-c", CMD)
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error while get container init: ", err)
		return "", err
	}
	initPID := string(out)
	return initPID, nil
}
