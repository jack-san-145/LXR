package handlers

import (
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"net/http"
	"time"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("comming inside create")
	var conBuilder models.ContainerBuilder
	err := json.NewDecoder(r.Body).Decode(&conBuilder.Container)
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

	//use seperate go routine to listen the quit signal for creation termination
	go func() {
		<-conBuilder.Quit
		conBuilder.Conn.Close()

	}()

	//check the container already exists or not
	exists := h.Helper.ContainerExists(conBuilder.Container.ContainerName)
	if exists {
		conn.Write([]byte("Container Already Exists"))
		return
	}

	conn.Write([]byte("creation started...\n"))

	//to check the image locally
	exists = h.Helper.CheckImageLocally(conBuilder.Container.Image)
	if !exists {
		err := h.Helper.PullImage(&conBuilder)
		if err != nil {
			conBuilder.Quit <- struct{}{} //pass quit signal to terminal container creation process
		}
	}

	//if not already exists creats new container environment (only setup containers rootfs)
	err = h.Helper.RootfsSetup(&conBuilder)
	if err != nil {
		log.Println("Error during RootfsSetup: ", err)
		conBuilder.Quit <- struct{}{} //pass quit signal to terminal container creation process
		return
	}

	//setup new container with rootfs
	err = h.Helper.ContainerSetup(conBuilder.Container)
	if err != nil {
		conBuilder.Quit <- struct{}{} //pass quit signal to terminal container creation process
		return
	}
	h.Helper.SetContainerActive(conBuilder.Container) //set new container to active state

	time.Sleep(time.Second * 3) //wait 3sec to complete container setup

	//setup networking for container
	err = h.Helper.SetupContainerNetworking(&conBuilder)
	if err != nil {
		conBuilder.Quit <- struct{}{}
		return
	}

	//create cgroups(memory,cpu,process)limits for new container
	err = h.Helper.CreateCgroup(conBuilder.Container)
	if err != nil {
		conBuilder.Quit <- struct{}{}
		return
	}

}
