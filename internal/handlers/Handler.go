package handlers

import "lxr-d/internal/models"

type LXRHandler struct {
	ContainerManager *models.ContainerManager
}

func NewHandler() *LXRHandler {
	return &LXRHandler{
		ContainerManager: models.NewContainerManager(),
	}
}
