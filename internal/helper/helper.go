package helper

import (
	"lxr-d/internal/ip"
	"lxr-d/internal/models"
)

type Helper struct {
	ContainerManager *models.ContainerManager
	NetworkConfig    *ip.NetworkConfig
}

func NewHelper(networkConfig *ip.NetworkConfig) *Helper {
	return &Helper{
		ContainerManager: models.NewContainerManager(),
		NetworkConfig:    ip.NewNetwork(*networkConfig),
	}
}
