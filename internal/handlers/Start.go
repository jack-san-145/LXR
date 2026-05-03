package handlers

import (
	"lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	con_name := r.URL.Query().Get("container_name")

	if container, ok := h.Helper.ContainerManager.AllContainers[con_name]; ok {
		err := h.Helper.ContainerSetup(container)
		if err != nil {
			response.WriteJson(w, map[string]bool{"active": false})
			return
		}
		h.Helper.ContainerManager.ActiveContainers[container.ContainerName] = &container.PID
		response.WriteJson(w, map[string]bool{"active": true})
	}
}
