package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/fatih/color"

	serverHandler "github.com/Bearaujus/minecraft-server-api/internal/handler/server"
	serverResource "github.com/Bearaujus/minecraft-server-api/internal/resource/server"
)

func main() {
	var serverResource = serverResource.NewServerResource()
	var serverHandler = serverHandler.NewServerHandler(serverResource)
	var router = NewRouter(serverHandler)

	address := "localhost:25001"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	fmt.Printf("service running at %v\n", color.YellowString(address))
	if err := http.Serve(listener, router); err != nil {
		return
	}
}
