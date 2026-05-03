package handlers

import (
	"encoding/json"
	"log"
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) PullImageHandler(w http.ResponseWriter, r *http.Request) {

	type img struct {
		ImageName string `json:"img_name"`
	}

	var image img
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		log.Println("Image Pull Error: ", err)
		return
	}
	log.Println("image from image handler: ", image.ImageName)

	exists, err := h.Helper.PullImage(image.ImageName)
	if exists {
		response.WriteJson(w, map[string]string{"status": "Image already exists locally"})
		return
	}
	if err != nil {
		response.WriteJson(w, map[string]string{"status": "Error in Image Pull"})
		return
	}
	response.WriteJson(w, map[string]string{"status": "Image Pulled Successfully"})

}
