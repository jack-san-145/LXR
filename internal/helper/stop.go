package helper

import (
	"log"
	"os"
	"strconv"
)

func (h *Helper) StopProcess(name string) (bool, error) {

	pid, exists := h.GetContainerPid(name)
	if !exists {
		log.Println("Process does not exists")
		return exists, nil
	}
	pid_int, _ := strconv.Atoi(pid)
	ps, err := os.FindProcess(pid_int)

	if err != nil {
		return exists, err
	}

	err = ps.Kill()
	if err == nil {
		delete(h.ContainerManager.ActiveContainers, name)
	}

	return exists, err
}
