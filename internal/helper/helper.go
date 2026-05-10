package helper

import (
	"lxr-d/internal/ip"
	"lxr-d/internal/models"
)

type Helper struct {
	ContainerManager *models.ContainerManager
	IpStack          *ip.IpStack
}

func NewHelper(stack *ip.IpStack) *Helper {
	return &Helper{
		ContainerManager: models.NewContainerManager(),
		IpStack:          ip.NewIpStack(*stack),
	}
}
