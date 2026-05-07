package models

type CreationResponse struct {
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	IsCreated     bool   `json:"is_created"`
	AlreadyExists bool   `json:"already_exists"`
}

type StartResponse struct {
	AlreadyActive bool `json:"already_active"`
	Activated     bool `json:"activated"`
	Failed        bool `json:"failed"`
	DoesNotExists bool `json:"exists"`
}

type StopResponse struct {
	Exists  bool `json:"container_exists"`
	Stopped bool `json:"container_Stoped"`
}
