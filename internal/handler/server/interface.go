package server

import "net/http"

type ServerHandlerItf interface {
	GetAllServerHandler(http.ResponseWriter, *http.Request) error
	CreateServerHandler(http.ResponseWriter, *http.Request) error
	DeleteServerHandler(http.ResponseWriter, *http.Request) error
	AgreeEulaServerHandler(http.ResponseWriter, *http.Request) error
	StartServerHandler(http.ResponseWriter, *http.Request) error
	StopServerHandler(http.ResponseWriter, *http.Request) error
	GetServerConsoleHandler(http.ResponseWriter, *http.Request) error
	AddServerConsoleHandler(http.ResponseWriter, *http.Request) error
}
