package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"lxr-d/internal/models"
	"net/http"
	"time"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {

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
	_, exists := h.Helper.ContainerExists(conBuilder.Container.ContainerName)
	if exists {
		conBuilder.Conn.Write([]byte("Container Already Exists\n"))
		conBuilder.Quit <- struct{}{}
		return
	}

	conBuilder.Conn.Write([]byte("creation started...\n"))

	//to check the image locally
	conBuilder.Conn.Write([]byte("[+] Find Image locally in LXR-registry...\n"))
	exists = h.Helper.CheckImageLocally(conBuilder.Container.Image)
	if !exists {
		err := h.Helper.PullImage(&conBuilder)
		if err != nil {
			conBuilder.Conn.Write([]byte("Image Pull Failed\n"))
			conBuilder.Quit <- struct{}{} //pass quit signal to terminal container creation process
		}

		//install dependency modules inside image rootfs
		err = h.Helper.InstallDependencies(&conBuilder)
		if err != nil {
			conBuilder.Conn.Write([]byte("Container dependency installation failed..\n"))
			conBuilder.Quit <- struct{}{}
			go h.Helper.KillContainer(conBuilder.Container) //kill failed container
			return
		}

	}

	//if not already exists creats new container environment (only setup containers rootfs)
	conBuilder.Conn.Write([]byte("\n[+] Setting up container rootfs...\n"))
	err = h.Helper.RootfsSetup(&conBuilder)
	if err != nil {
		conBuilder.Conn.Write([]byte("Rootfs setup failed..\n"))
		conBuilder.Quit <- struct{}{}                   //pass quit signal to terminal container creation process
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	//setup new container with rootfs
	conBuilder.Conn.Write([]byte("\n[+] Building container environment with rootfs..\n"))
	err = h.Helper.ContainerSetup(conBuilder.Container)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container setup failed..\n"))
		conBuilder.Quit <- struct{}{}                   //pass quit signal to terminal container creation process
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	//add containers to Allcontainers in containerManager
	h.Helper.AddContainer(conBuilder.Container)

	h.Helper.SetContainerActive(conBuilder.Container.ContainerName) //set new container to active state

	time.Sleep(time.Second * 3) //wait 3sec to complete container setup

	//setup networking for container
	conBuilder.Conn.Write([]byte("[+] Setting up container networking..\n"))
	err = h.Helper.SetupContainerNetworking(&conBuilder)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container networking failed..\n"))
		conBuilder.Quit <- struct{}{}
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	//create cgroups(memory,cpu,process)limits for new container
	conBuilder.Conn.Write([]byte("[+] Setting up container resources limit...\n"))
	err = h.Helper.CreateCgroup(conBuilder.Container)
	if err != nil {
		conBuilder.Conn.Write([]byte("Container cgroups creation failed..\n"))
		conBuilder.Quit <- struct{}{}
		go h.Helper.KillContainer(conBuilder.Container) //kill failed container
		return
	}

	// //freeze container immediately after container creation
	// h.Helper.FreezeContainer(conBuilder.Container.ContainerName)

	conBuilder.Container.Ports = append(conBuilder.Container.Ports, 9000) //add exposed code-server port(9000)

	conBuilder.Conn.Write([]byte("[+] code-server activated at port 9000 ✔\n"))
	containerDetails := fmt.Sprintf("\nCONTAINER ID: %v              CONTAINER NAME: %v", conBuilder.Container.ContainerId, conBuilder.Container.ContainerName)
	conBuilder.Conn.Write([]byte(containerDetails + "\n"))

	conBuilder.Conn.Write([]byte("\nContainer Created Successfully...\n"))
	conBuilder.Quit <- struct{}{}
}
