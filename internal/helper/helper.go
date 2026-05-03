package helper

import "lxr-d/internal/models"

type Helper struct {
	ContainerManager *models.ContainerManager
}

func NewHelper() *Helper {
	return &Helper{
		ContainerManager: models.NewContainerManager(),
	}
}
