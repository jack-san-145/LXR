package helper

import (
	"lxr-d/internal/models"
	"os/exec"
	"strconv"
	"strings"
)

func (h *Helper) GetContainerPid(name string) (string, bool) {

	h.ContainerManager.Mu.RLock()
	defer h.ContainerManager.Mu.RUnlock()

	pid, ok := h.ContainerManager.ActiveContainers[name]

	if pid != 0 {
		return strconv.Itoa(pid), ok
	}
	return "", false
}

// to get child PID from Parent PID
func (h *Helper) GetChildPID(ParentPID int) (string, error) {

	cmd := exec.Command(
		"pgrep",
		"-P",
		strconv.Itoa(ParentPID),
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	//returns only first child of parent
	return strings.TrimSpace(string(out)), nil
}

// check container exists or not
func (h *Helper) ContainerExists(name string) bool {

	h.ContainerManager.Mu.RLock()
	defer h.ContainerManager.Mu.RUnlock()

	_, ok := h.ContainerManager.AllContainers[name]
	return ok
}

// check container currently active or not
func (h *Helper) IsContainerActive(name string) bool {

	h.ContainerManager.Mu.RLock()
	defer h.ContainerManager.Mu.RUnlock()

	_, ok := h.ContainerManager.ActiveContainers[name]
	return ok
}

// add newly created container to allContainers
func (h *Helper) AddContainer(con *models.Container) {

	h.ContainerManager.Mu.Lock()
	defer h.ContainerManager.Mu.Unlock()

	h.ContainerManager.AllContainers[con.ContainerName] = con
}

// add container to ActiveContainers
func (h *Helper) SetContainerActive(containerName string) {

	h.ContainerManager.Mu.Lock()
	defer h.ContainerManager.Mu.Unlock()

	PID := h.ContainerManager.AllContainers[containerName].PID
	h.ContainerManager.ActiveContainers[containerName] = PID
}

// remove container from ActiveContainers
func (h *Helper) SetContainerDeactive(containerName string) {

	h.ContainerManager.Mu.Lock()
	defer h.ContainerManager.Mu.Unlock()

	delete(h.ContainerManager.ActiveContainers, containerName)
}

// remove container to AllContainers
func (h *Helper) RemoveContainer(containerName string) {

	h.ContainerManager.Mu.Lock()
	defer h.ContainerManager.Mu.Unlock()

	delete(h.ContainerManager.AllContainers, containerName)
}
