package main

import (
	"encoding/json"
	"net/http"

	serverHandler "github.com/Bearaujus/minecraft-server-api/internal/handler/server"
	"github.com/Bearaujus/minecraft-server-api/internal/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type httpHandler func(w http.ResponseWriter, r *http.Request) error

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		res, _ := json.Marshal(model.Response{
			Header: model.ResponseHeader{
				ProcessTime: w.Header().Get("time_elapsed"),
				IsSuccess:   false,
				Messages:    err.Error(),
			},
			Data: nil,
		})
		w.Header().Del("time_elapsed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}
}

func NewRouter(sh serverHandler.ServerHandlerItf) *chi.Mux {
	router := chi.NewRouter()
	router.MethodNotAllowed(http.NotFound)
	router.Use(middleware.Logger)

	// get all servers data
	router.Method(http.MethodGet, "/servers", httpHandler(sh.GetAllServerHandler))
	// create new server
	router.Method(http.MethodPost, "/servers/create", httpHandler(sh.CreateServerHandler))
	// delete server
	router.Method(http.MethodDelete, "/server/{id}/delete", httpHandler(sh.DeleteServerHandler))
	// agree eula
	router.Method(http.MethodPatch, "/server/{id}/agree-eula", httpHandler(sh.AgreeEulaServerHandler))
	// start server
	router.Method(http.MethodPatch, "/server/{id}/start", httpHandler(sh.StartServerHandler))
	// stop server
	router.Method(http.MethodPatch, "/server/{id}/stop", httpHandler(sh.StopServerHandler))
	// get current server console status
	router.Method(http.MethodGet, "/server/{id}/console", httpHandler(sh.GetServerConsoleHandler))
	// add command to console
	router.Method(http.MethodPost, "/server/{id}/console/execute", httpHandler(sh.AddServerConsoleHandler))

	return router
}
