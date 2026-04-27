package models

type Container struct {
	Image         string `json:"image_name"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_id"`
	Active        bool   `json:"active"`
	Port          int    `json:"port"`
	IpAddress     string `json:"ip_address"`
	Bridge        string `json:"bridge"`
}

type ContainerManager struct {
	AllContainers    map[string]*Container
	ActiveContainers map[string]*Container
}

func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		AllContainers:    map[string]*Container{},
		ActiveContainers: map[string]*Container{},
	}
}
