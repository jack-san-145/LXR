package helper

import (
	"fmt"
	"log"

	"os/exec"
	"strconv"
)

// to find equivalent bytes for given mb
func (h *Helper) MB(m int) string {
	return strconv.Itoa(m * 1024 * 1024) //mb -> kb -> byte
}

// to find percent for given integer
func (h *Helper) CPU(percent int) string {
	quota := percent * 1000
	return fmt.Sprintf("%v 100000", quota)
}

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
