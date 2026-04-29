package main

import (
	"context"
	"log"
	"lxr-d/internal/handlers"
	"net"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {

	lxr_sock := "/var/run/lxr.sock"
	err := os.Remove(lxr_sock)
	if err != nil {
		log.Println("Old Sock remove Error: ", err)
	}

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

	LXRHandler := handlers.NewHandler()
	router := NewRouter(LXRHandler)
	go runServer(router, listener)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	<-ctx.Done()
	LXRHandler.Helper.SaveContainerState()
	defer stop()

}

func runServer(router *chi.Mux, listener net.Listener) {
	//start the server
	log.Println("Server Listening ....")
	err := http.Serve(listener, router)
	if err != nil {
		log.Fatal("Server Failed to start: ", err)
	}
}
