package helper

import (
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"os"
)

func (h *Helper) SaveContainerState() {

	h.ContainerManager.Mu.RLock()
	defer h.ContainerManager.Mu.RUnlock()

	backup := models.BackupContainerState{
		ContainerManager: models.BackupContainerManager{
			AllContainers:    h.ContainerManager.AllContainers,
			ActiveContainers: h.ContainerManager.ActiveContainers,
		},
		NetworkConfig: *h.NetworkConfig,
	}

	file, err := os.Create("/home/LXR/Container-state.json")
	if err != nil {
		log.Println("JSON creation error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(&backup); err != nil {
		log.Println("JSON encode error:", err)
	}
}

func (h *Helper) BackupContainerState() {

	file, err := os.Open("/home/LXR/Container-state.json")
	if err != nil {
		log.Println("JSON extraction error:", err)
		return
	}
	defer file.Close()

	var backup models.BackupContainerState

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&backup); err != nil {
		log.Println("JSON decode error:", err)
		return
	}

	//check if states are nill ,if true then declare empty map
	if backup.ContainerManager.AllContainers == nil {
		backup.ContainerManager.AllContainers = make(map[string]*models.Container)
	}

	if backup.ContainerManager.ActiveContainers == nil {
		backup.ContainerManager.ActiveContainers = make(map[string]int)
	}

	h.ContainerManager.Mu.Lock()
	defer h.ContainerManager.Mu.Unlock()

	//assigns existing container state values to current container state
	h.ContainerManager.AllContainers = backup.ContainerManager.AllContainers
	h.ContainerManager.ActiveContainers = backup.ContainerManager.ActiveContainers

	//store backupState ipStack to current runtime containerstate
	h.NetworkConfig.IPStack = backup.NetworkConfig.IPStack

	if backup.NetworkConfig.LastUsedIP != "" {
		h.NetworkConfig.LastUsedIP = backup.NetworkConfig.LastUsedIP
	}

	if backup.NetworkConfig.HostUsed != 0 {
		h.NetworkConfig.HostUsed = backup.NetworkConfig.HostUsed
	}
}
