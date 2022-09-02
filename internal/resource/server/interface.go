package server

import "github.com/Bearaujus/minecraft-server-api/internal/model"

type ServerResourceItf interface {
	GetAllServerResource() (map[string]*model.Server, error)
	CreateServerResource() (string, error)
	DeleteServerResource(string) error
	AgreeEulaServerResource(string) error
	StartServerResource(string, int, int, string) error
	StopServerResource(string) error
	GetServerConsoleResource(string) ([]byte, error)
	AddServerConsoleResource(string, string) error
}
