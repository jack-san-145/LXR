package helper

import (
	"lxr-d/internal/models"
	"os/exec"
	"strconv"
	"strings"
)

func (h *Helper) GetContainerPid(name string) (string, bool) {

	pid, ok := h.ContainerManager.ActiveContainers[name]

	if pid != nil {
		return strconv.Itoa(*pid), ok
	}
	return "", false
}

// to get container init PID
func (h *Helper) GetContainerInitPid(unsharePID int) (string, error) {

	cmd := exec.Command(
		"pgrep",
		"-P",
		strconv.Itoa(unsharePID),
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
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

// add newly created container to allContainers
func (h *Helper) AddContainer(con *models.Container) {
	h.ContainerManager.AllContainers[con.ContainerName] = con
}

// add container to ActiveContainers
func (h *Helper) SetContainerActive(con *models.Container) {
	h.ContainerManager.ActiveContainers[con.ContainerName] = &con.PID
}
