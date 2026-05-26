package helper

import "lxr-d/internal/models"

// returns only running containers
func (h *Helper) GetActiveContainers() []models.PsContainer {

	var containers []models.PsContainer

	for name := range h.ContainerManager.ActiveContainers {

		con, ok := h.ContainerManager.AllContainers[name]
		if !ok {
			continue
		}

		containers = append(containers, models.PsContainer{
			ContainerID:   con.ContainerId,
			ContainerName: con.ContainerName,
			Image:         con.Image,
			PID:           con.PID,
			Status:        "running",
			IPAddress:     con.IpAddress,
			Port:          con.Port,
		})
	}

	return containers
}
