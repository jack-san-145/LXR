package helper

import "lxr-d/internal/models"

func (h *Helper) GetContainerByName(name string) {

}

type Helper struct {
	ContainerManager *models.ContainerManager
}

func NewHelper() *Helper {
	return &Helper{
		ContainerManager: models.NewContainerManager(),
	}
}
