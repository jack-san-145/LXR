package helper

import (
	"log"
	"os"
)

func checkImageLocally(image string) bool {
	path := "/home/LXR/LXR-registry/" + image
	info, err := os.Stat(path)

	if err != nil {
		log.Println("Error checking folder: ", err)
		return false
	}

	if info.IsDir() {
		return true
	}
	return false
}
