package helper

import (
	"strconv"
)

func (h *Helper) GetContainerPid(name string) (string, bool) {

	pid, ok := h.ContainerManager.ActiveContainers[name]

	return strconv.Itoa(*pid), ok
}
