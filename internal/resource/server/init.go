package server

import (
	"errors"
	"fmt"

	"github.com/Bearaujus/minecraft-server-api/internal/model"
	"github.com/Bearaujus/minecraft-server-api/pkg"
)

type serverResource struct {
	serverdata map[string]*model.Server
}

func NewServerResource() ServerResourceItf {
	var res = &serverResource{
		serverdata: make(map[string]*model.Server),
	}

	pkg.ValidateDir(true, model.DIR_SERVER)
	var modelServerID, _ = pkg.GetListFolderFromDir(model.DIR_SERVER)
	for _, id := range modelServerID {
		res.serverdata[id] = nil
	}

	return res
}

func (sr *serverResource) addServer(id string) error {
	var _, ok = sr.serverdata[id]
	if ok {
		return errors.New("server already exist")
	}
	sr.serverdata[id] = nil

	return nil
}

func (sr *serverResource) getServer(id string) (*model.Server, error) {
	var res, ok = sr.serverdata[id]
	if !ok {
		return nil, errors.New("server not exist")
	}

	if res != nil {
		// kill server if pipe is broken
		if _, err := fmt.Fprint(*res.StdinPipe, ""); err != nil {
			res.Cmd.Process.Kill()
			sr.serverdata[id] = nil
		}
	}

	return res, nil
}

func (sr *serverResource) setServer(id string, server *model.Server) error {
	sr.serverdata[id] = server

	return nil
}
