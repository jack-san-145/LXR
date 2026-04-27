package models

import ()

type CreationResponse struct {
	IsCreated     bool   `json:"is_created"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_id"`
}
