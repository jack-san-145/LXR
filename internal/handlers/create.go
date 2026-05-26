package handlers

import (
	"bufio"
	"encoding/json"
	"log"
	"lxr-d/internal/models"
	"net/http"
	"time"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("comming inside create")
	var (
		conBuilder models.ContainerBuilder
		buf        *bufio.ReadWriter
	)

	conBuilder.Quit = make(chan struct{}) //intialize quit channel

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

	//use seperate go routine to listen the quit signal for creation termination
	go func() {
		<-conBuilder.Quit
		conBuilder.Conn.Close()

	}()

	//check the container already exists or not
	exists := h.Helper.ContainerExists(conBuilder.Container.ContainerName)
	if exists {
		conBuilder.Conn.Write([]byte("Container Already Exists\n"))
		return
	}

	conBuilder.Conn.Write([]byte("creation started...\n"))

	//to check the image locally
	conBuilder.Conn.Write([]byte("[+]Find Image locally in LXR-registry...\n"))
	exists = h.Helper.CheckImageLocally(conBuilder.Container.Image)
	if !exists {
		err := h.Helper.PullImage(&conBuilder)
		if err != nil {
			conBuilder.Conn.Write([]byte("Image Pull Failed\n"))
			conBuilder.Quit <- struct{}{} //pass quit signal to terminal container creation process
		}
	}

	//if not already exists creats new container environment (only setup containers rootfs)
	conBuilder.Conn.Write([]byte("[+]Setting up container rootfs...\n"))
	err = h.Helper.RootfsSetup(&conBuilder)
	if err != nil {
		conBuilder.Conn.Write([]byte("Rootfs setup failed..\n"))
		conBuilder.Quit <- struct{}{}                   //pass quit signal to terminal container creation process
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	//setup new container with rootfs
	conBuilder.Conn.Write([]byte("[+]Building container environment with rootfs..\n"))
	err = h.Helper.ContainerSetup(conBuilder.Container)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container setup failed..\n"))
		conBuilder.Quit <- struct{}{}                   //pass quit signal to terminal container creation process
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}
	h.Helper.SetContainerActive(conBuilder.Container) //set new container to active state

	time.Sleep(time.Second * 3) //wait 3sec to complete container setup

	//setup networking for container
	conBuilder.Conn.Write([]byte("\n[+]Setting up container networking..\n"))
	err = h.Helper.SetupContainerNetworking(&conBuilder)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container networking failed..\n"))
		conBuilder.Quit <- struct{}{}
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	//create cgroups(memory,cpu,process)limits for new container
	conBuilder.Conn.Write([]byte("\n[+]Setting up container resources limit...\n"))
	err = h.Helper.CreateCgroup(conBuilder.Container)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container cgroups creation failed..\n"))
		conBuilder.Quit <- struct{}{}
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	conBuilder.Conn.Write([]byte("\nContainer Created Successfully...\n"))
	conBuilder.Quit <- struct{}{}

	//add containers to Allcontainers in containerManager
	h.Helper.SetContainerDeactive(conBuilder.Container)
	h.Helper.AddContainer(conBuilder.Container)
}
