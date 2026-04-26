package models

import ()

type NewContainer struct {
	Image         string `json:"image_name"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_id"`
	HostName      string `json:"host_name"`
}

type CreationResponse struct {
	IsCreated     bool   `json:"is_created"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_id"`
}
