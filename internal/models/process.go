package models

type PsContainer struct {
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	Image         string `json:"image"`
	PID           int    `json:"pid"`
	Status        string `json:"status"`
	IPAddress     string `json:"ip_address"`
	Ports         []int  `json:"ports"`
}

type PsResponse struct {
	Containers []PsContainer `json:"containers"`
}
