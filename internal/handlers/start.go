package handlers

import (
	"bufio"
	"log"
	"lxr-d/internal/models"
	// "lxr-d/internal/response"
	"net/http"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	// con_name := r.URL.Query().Get("container_name")
	var (
		conBuilder models.ContainerBuilder
		buf        *bufio.ReadWriter
		err        error
	)

	//hijack the http connnection and stream i/o in real-time over uds
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack not supported", http.StatusInternalServerError)
		return
	}
	conBuilder.Conn, buf, err = hijacker.Hijack()
	if err != nil {
		log.Println("Hijack error:", err)
		return
	}

	//write http success response manually
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Type: text/plain\r\n")
	buf.WriteString("\r\n")
	buf.Flush()

	conBuilder.Conn.Write([]byte("Connection Hijacked"))

}
