package server

import (
	serverResource "github.com/Bearaujus/minecraft-server-api/internal/resource/server"
)

type serverHandler struct {
	Resource serverResource.ServerResourceItf
}

func NewServerHandler(resource serverResource.ServerResourceItf) ServerHandlerItf {
	return &serverHandler{
		Resource: resource,
	}
}
