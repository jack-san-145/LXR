package helper

import (
	"os"
	"strconv"
)

func (h *Helper) StopContainer(name string) (bool, error) {

	pid, exists := h.GetContainerPid(name)
	if !exists {
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
