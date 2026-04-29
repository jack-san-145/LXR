package helper

import (
	"encoding/json"
	"log"
	"os"
)

func (h *Helper) SaveContainerState() {

	log.Println("calling container state")
	file, err := os.Create("/home/jack/LXR-data/Container-state.json")
	if err != nil {
		log.Println("Json creation error: ", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	active_containers := h.ContainerManager.ActiveContainers
	all_containers := h.ContainerManager.AllContainers

	state_json := map[string]any{"Active_containers": active_containers, "all_containers": all_containers}

	log.Println("state_json ", state_json)
	err = encoder.Encode(state_json)
	if err != nil {
		log.Println("Json encode error: ", err)

	}
}
