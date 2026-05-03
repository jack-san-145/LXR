package handlers

import (
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) PullImageHandler(w http.ResponseWriter, r *http.Request) {

	img_name := r.URL.Query().Get("image_name")
	exists, err := h.Helper.PullImage(img_name)
	if exists {
		response.WriteJson(w, map[string]string{"status": "Image already exists locally"})
		return
	}
	if err != nil {
		response.WriteJson(w, map[string]string{"status": "Error in Image Pull"})
		return
	}
	response.WriteJson(w, map[string]string{"status": "Image Created Successfully"})

}
