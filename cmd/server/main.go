package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func main() {

	lxr_sock := "/var/run/lxr.sock"
	err := os.Remove(lxr_sock)
	if err != nil {
		log.Println("Old Sock remove Error: ", err)
	}

	r := chi.NewRouter()
	r.Get("/ping", PingHanlder)

	listener, err := net.Listen("unix", lxr_sock)
	if err != nil {
		log.Fatal("Listener Failed: ", err)
	}

	//changing sock path permissions
	group, err := user.LookupGroup("lxr")
	if err != nil {
		log.Println("Group not found: ", err)
	}
	group_id, _ := strconv.Atoi(group.Gid)

	err = os.Chown(lxr_sock, 0, group_id)
	if err != nil {
		log.Println("chown Error: ", err)
	}
	err = os.Chmod(lxr_sock, 0660)
	if err != nil {
		log.Println("chmod Error: ", err)
	}

	//start the server
	log.Println("Server Listening ....")
	err = http.Serve(listener, r)
	if err != nil {
		log.Fatal("Server Failed to start: ", err)
	}
}

func PingHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("calling")
	w.Write([]byte("Pong\n"))
}
