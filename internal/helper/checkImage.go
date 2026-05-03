package helper

import (
	"log"
	"os"
)

func (h *Helper) CheckImageLocally(image string) bool {
	path := "/home/LXR/LXR-registry/" + image
	log.Println("image: ", image)
	info, err := os.Stat(path)
	log.Println("path: ", path)

	if err != nil {
		log.Println("Error checking folder: ", err)

		if os.IsNotExist(err) {
			log.Println("Folder does NOT exist")
		} else {
			log.Println("Error:", err)
		}
		return false
	}

	if info.IsDir() {
		return true

	}
	return false
}
