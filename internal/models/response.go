package models

import ()

type CreationResponse struct {
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	IsCreated     bool   `json:"is_created"`
}
