package handlers

import "net/http"

func (h *LXRHandler) PingHanlder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong\n"))
}
