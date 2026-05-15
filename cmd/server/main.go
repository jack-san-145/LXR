package main

import (
	"context"
	"log"
	"lxr-d/internal/app"
	"lxr-d/internal/ip"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	//to load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env")
	}

	cidr, _ := strconv.Atoi(os.Getenv("CIDR"))
	totalUsableHost, _ := strconv.Atoi(os.Getenv("TOTAL_USABLE_HOST"))

	//create IpStack with own lxr network
	ipStack := ip.IpStack{
		Network:         os.Getenv("NETWORK"),
		Cidr:            cidr,
		BridgeName:      "lxr0",
		BridgeIp:        os.Getenv("BRIDGE_IP"),
		IpStartRange:    os.Getenv("IP_START_RANGE"),
		IpEndRange:      os.Getenv("IP_END_RANGE"),
		NetworkAddr:     os.Getenv("NETWORK_ADDR"),
		BroadCastAddr:   os.Getenv("BROADCAST_ADDR"),
		TotalUsableHost: totalUsableHost,
	}

	//create a runtime named Lxr with its own ipStack
	Lxr := app.NewApp(&ipStack)

	//creates a chi router with LXR handler
	router := NewRouter(Lxr.Handler)

	listener := Lxr.Helper.InitDaemon() //start the daemon initialization and return listener for unix sock connection

	Lxr.Helper.BackupContainerState() //backup existing container state

	go runServer(router, listener) //start the go server in seperate go routine

	//context to listen the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	//when interrupt occurs save container state ,then stop the daemon
	Lxr.Helper.SaveContainerState()
}

func runServer(router *chi.Mux, listener net.Listener) {
	//start the server
	log.Println("Server Listening ....")
	err := http.Serve(listener, router)
	if err != nil {
		log.Fatal("Server Failed to start: ", err)
	}
}
