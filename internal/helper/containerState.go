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
	err = encoder.Encode(h.ContainerManager)
	if err != nil {
		log.Println("Json encode error: ", err)

	}
}

func (h *Helper) BackupContainerState() {

	file, err := os.Open("/home/jack/LXR-data/Container-state.json")
	if err != nil {
		log.Println("Json extraction error: ", err)
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(h.ContainerManager)
	if err != nil {
		log.Println("JSON decode error: ", err)
	}

}
