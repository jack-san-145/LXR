package main

import (
	"context"
	"log"
	"lxr-d/internal/app"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {

	//create a runtime named Lxr
	Lxr := app.NewApp()

	//creates a chi router
	router := NewRouter(Lxr.Handler)

	listener := Lxr.Helper.InitDaemon() //start the daemon initialization

	Lxr.Helper.BackupContainerState() //backup existing container state

	go runServer(router, listener) //start the go server in seperate go routine

	//context to listen the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()

	//when interrupt occurs save container state ,then stop the daemon
	Lxr.Helper.SaveContainerState()
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
