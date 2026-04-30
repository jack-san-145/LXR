package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"lxr-d/internal/handlers"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	LXRHandler := handlers.NewHandler()
	router := NewRouter(LXRHandler)

	listener := LXRHandler.Helper.InitDaemon()
	LXRHandler.Helper.BackupContainerState()

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
