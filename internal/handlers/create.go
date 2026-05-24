package handlers

import (
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"net/http"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("comming inside create")
	var con models.Container
	err := json.NewDecoder(r.Body).Decode(&con)
	if err != nil {
		log.Println("Creation Error: ", err)
		return
	}

	//hijack the http connnection and stream i/o in real-time over uds
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack not supported", http.StatusInternalServerError)
		return
	}
	conn, buf, err := hijacker.Hijack()
	if err != nil {
		log.Println("Hijack error:", err)
		return
	}
	defer conn.Close()

	//write http success response manually
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Type: text/plain\r\n")
	buf.WriteString("\r\n")
	buf.Flush()

	conn.Write([]byte("creation started...\n"))

}
