package helper

import (
	"strconv"
)

func (h *Helper) GetContainerPid(name string) (string, bool) {

	pid, ok := h.ContainerManager.ActiveContainers[name]

	if pid != nil {
		return strconv.Itoa(*pid), ok
	}
	return "", false
}

// check container exists or not
func (h *Helper) ContainerExists(name string) bool {

	_, ok := h.ContainerManager.AllContainers[name]
	return ok
}

// check container currently active or not
func (h *Helper) ContainerActive(name string) bool {
	_, ok := h.ContainerManager.ActiveContainers[name]
	return ok
}
