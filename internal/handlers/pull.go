package handlers

import (
	"bufio"
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"net/http"
)

func (h *Handler) PullImageHandler(w http.ResponseWriter, r *http.Request) {

	type img struct {
		ImageName string `json:"img_name"`
	}

	var (
		image   img
		builder models.ContainerBuilder
		buf     *bufio.ReadWriter
	)

	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		log.Println("Image Pull Error: ", err)
		return
	}

	builder.Container.Image = image.ImageName
	//hijack the http connnection and stream i/o in real-time over uds
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack not supported", http.StatusInternalServerError)
		return
	}
	builder.Conn, buf, err = hijacker.Hijack()
	if err != nil {
		log.Println("Hijack error:", err)
		return
	}
	defer builder.Conn.Close()

	//write http success response manually
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Type: text/plain\r\n")
	buf.WriteString("\r\n")
	buf.Flush()

	//use seperate go routine to listen the quit signal for creation termination
	go func() {
		<-builder.Quit
		builder.Conn.Close()

	}()

	exists := h.Helper.CheckImageLocally(builder.Container.Image)
	if exists {
		builder.Conn.Write([]byte("Image already exists locally..\n"))
		builder.Quit <- struct{}{}
		return
	}

	err = h.Helper.PullImage(&builder)

	if err != nil {
		builder.Conn.Write([]byte("Error in Image Pull\n"))
		builder.Quit <- struct{}{}
		return
	}
	builder.Conn.Write([]byte("Image Pulled Successfully..\n"))
	builder.Quit <- struct{}{}

}
