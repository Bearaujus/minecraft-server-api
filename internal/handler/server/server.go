package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Bearaujus/minecraft-server-api/internal/model"
	"github.com/Bearaujus/minecraft-server-api/pkg"

	"github.com/go-chi/chi"
)

func (sh *serverHandler) GetAllServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	modelServer, err := sh.Resource.GetAllServerResource()
	if err != nil {
		return err
	}

	outputRes := make([]model.GetAllServerResponse, 0)
	for k, v := range modelServer {
		resItem := model.GetAllServerResponse{
			ServerID: k,
		}

		if v == nil {
			resItem.Status = "stopped"
			resItem.LastError = resItem.GetLastError(k)
		} else if v.IsAttemptedToStop {
			resItem.Status = "stopping"
		} else if v.IsAttemptedToStart {
			resItem.Status = "starting"
		} else {
			resItem.Status = "running"
			resItem.Address = fmt.Sprintf("localhost:%v", v.Port)
			resItem.OnlineMode = resItem.IsRunningOnlineMode(k)
			resItem.WorldName = resItem.GetUsedWorldName(k)
		}

		outputRes = append(outputRes, resItem)
	}

	sort.Slice(outputRes, func(i, j int) bool {
		return outputRes[i].ServerID < outputRes[j].ServerID
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: outputRes,
	})
}

func (sh *serverHandler) CreateServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	res, err := sh.Resource.CreateServerResource()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: res,
	})
}

func (sh *serverHandler) DeleteServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	err := sh.Resource.DeleteServerResource(id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: "server successfully deleted",
	})
}

func (sh *serverHandler) AgreeEulaServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	err := sh.Resource.AgreeEulaServerResource(id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: "eula set to agree",
	})
}

func (sh *serverHandler) StartServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	// parse ram
	sRamGB := r.FormValue("ram_gb")
	if sRamGB == "" {
		return errors.New("ram_gb is required")
	}
	ramGB, err := strconv.Atoi(sRamGB)
	if err != nil {
		return err
	}
	if ramGB <= 0 {
		return errors.New("ram_gb cannot <= 0")
	}

	// parse port
	sPort := r.FormValue("port")
	if sPort == "" {
		return errors.New("port is required")
	}
	port, err := strconv.Atoi(sPort)
	if err != nil {
		return err
	}
	if port < 25000 {
		return errors.New("port cannot <= 25000")
	}
	if port > 30000 {
		return errors.New("port cannot >= 30000")
	}

	// parse world
	worldName := r.FormValue("world_name")

	if err := sh.Resource.StartServerResource(id, ramGB, port, worldName); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: "attempted to start",
	})
}

func (sh *serverHandler) StopServerHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	if err := sh.Resource.StopServerResource(id); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: "attempted to stop",
	})
}

func (sh *serverHandler) GetServerConsoleHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	res, err := sh.Resource.GetServerConsoleResource(id)
	if err != nil {
		return err
	}

	// parse limit
	sLimit := r.URL.Query().Get("limit")
	if sLimit != "" {
		limit, err := strconv.Atoi(sLimit)
		if err != nil {
			return err
		}
		if limit <= 0 {
			return errors.New("limit cannot <= 0")
		}

		tmpRes := strings.Split(string(res), "\n")
		if len(tmpRes) >= limit {
			tmpRes = tmpRes[len(tmpRes)-limit:]
		}

		res = []byte(strings.Join(tmpRes, "\n"))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return nil
}

func (sh *serverHandler) AddServerConsoleHandler(w http.ResponseWriter, r *http.Request) error {
	timer := pkg.StartNewTimer()
	defer func() {
		w.Header().Add("time_elapsed", timer.SinceStringInMS())
	}()

	// parse id
	id := chi.URLParam(r, "id")
	if id == "" {
		return errors.New("id is required")
	}

	// parse command
	command := r.FormValue("command")
	if command == "" {
		return errors.New("command is required")
	}

	if err := sh.Resource.AddServerConsoleResource(id, command); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(model.Response{
		Header: model.ResponseHeader{
			ProcessTime: timer.SinceStringInMS(),
			IsSuccess:   true,
			Messages:    nil,
		},
		Data: "command executed",
	})
}
