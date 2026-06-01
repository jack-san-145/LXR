package handlers

import (
	"bufio"
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"net/http"
	// "time"
)

func (h *Handler) StartHandler(w http.ResponseWriter, r *http.Request) {

	type conNameStruct struct {
		ContainerName string `json:"container_name"`
	}

	var (
		conBuilder models.ContainerBuilder
		buf        *bufio.ReadWriter
		conName    conNameStruct
	)

	//initialize conBuilder values
	conBuilder.Quit = make(chan struct{})

	err := json.NewDecoder(r.Body).Decode(&conName)
	if err != nil {
		log.Println("Start container Error: ", err)
		return
	}

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

	//use seperate go routine to listen the quit signal for worlflow termination
	go func() {
		<-conBuilder.Quit
		conBuilder.Conn.Close()

	}()

	//check whether the container is exists or not ,if exists make it active
	container, exists := h.Helper.ContainerExists(conName.ContainerName)
	if !exists {
		conBuilder.Conn.Write([]byte("Container Doesn't exists\n"))
		conBuilder.Quit <- struct{}{}
		return
	}

	conBuilder.Container = container

	//check whether the container is active or not
	active := h.Helper.IsContainerActive(conBuilder.Container.ContainerName)
	if active {
		conBuilder.Conn.Write([]byte("Container Already Running...\n"))
		conBuilder.Quit <- struct{}{}
		return
	}

	if container.Freezed {
		h.Helper.UnfreezeContainer(container.ContainerName)
		conBuilder.Conn.Write([]byte("Container Running...\n"))
		conBuilder.Quit <- struct{}{}
		return
	}

}
