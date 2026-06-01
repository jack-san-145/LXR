package models

import (
	"net"
	"sync"
)

type Container struct {
	PID           int    `json:"pid"`
	Image         string `json:"image_name"`
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	Bridge        string `json:"bridge"`
	IpAddress     string `json:"ip_address"`
	Ports         []int  `json:"ports"`
	Freezed       bool   `json:"freezed"`
	Active        bool   `json:"active"`
	ConVeth       string `json:"con_veth"`
	BrVeth        string `json:"br_veth"`
}

// container creater with client connection and channel
type ContainerBuilder struct {
	Container *Container
	Conn      net.Conn
	Quit      chan struct{}
}

type ContainerManager struct {
	AllContainers    map[string]*Container
	ActiveContainers map[string]int
	Mu               sync.RWMutex
}

func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		AllContainers:    map[string]*Container{},
		ActiveContainers: map[string]int{},
	}
}
