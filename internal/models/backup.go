package models

import "lxr-d/internal/ip"

type BackupContainerManager struct {
	AllContainers    map[string]*Container `json:"all_containers"`
	ActiveContainers map[string]int        `json:"active_containers"`
}

type BackupContainerState struct {
	ContainerManager BackupContainerManager `json:"container_manager"`
	NetworkConfig    ip.NetworkConfig       `json:"network_config"`
}
