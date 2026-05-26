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

// returns all containers
func (h *Helper) GetAllContainers() []models.PsContainer {

	var containers []models.PsContainer

	for name, con := range h.ContainerManager.AllContainers {

		status := "stopped"

		if _, ok := h.ContainerManager.ActiveContainers[name]; ok {
			status = "running"
		}

		containers = append(containers, models.PsContainer{
			ContainerID:   con.ContainerId,
			ContainerName: con.ContainerName,
			Image:         con.Image,
			PID:           con.PID,
			Status:        status,
			IPAddress:     con.IpAddress,
			Port:          con.Port,
		})
	}

	return containers
}
