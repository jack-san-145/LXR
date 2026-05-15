package helper

import (
	"log"
	"lxr-d/internal/models"
)

func (h *Helper) SetupContainerNetworking(con *models.Container) error {

	//use lxr defaulat bridge
	con.Bridge = h.IpStack.UseBridge()

	//allocate ip address for container
	con.IpAddress = h.IpStack.AllocateIp()

	log.Println("container ip: ", con.IpAddress)

	return nil
}
