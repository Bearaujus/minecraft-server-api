package model

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
)

var (
	DIR_JAR    = path.Join("file", "jar")
	DIR_SERVER = path.Join("file", "server")
)

type Server struct {
	ID    string
	Port  int
	RamGB int

	Cmd       *exec.Cmd
	StdinPipe *io.WriteCloser
	FileOut   *os.File

	IsAttemptedToStart bool
	IsAttemptedToStop  bool
}

type GetAllServerResponse struct {
	ServerID   string `json:"server_id"`
	Status     string `json:"status"`
	Address    string `json:"address,omitempty"`
	OnlineMode bool   `json:"online_mode,omitempty"`
	WorldName  string `json:"world_name,omitempty"`
	LastError  string `json:"last_error,omitempty"`
}

func (gasr *GetAllServerResponse) GetLastError(id string) string {
	if !gasr.isEulaAccepted(id) {
		return "need to accept eula"
	}

	if gasr.isPortAlreadyUsed(id) {
		return "fail to bind port"
	}

	return ""
}

func (gasr *GetAllServerResponse) isEulaAccepted(id string) bool {
	data, err := ioutil.ReadFile(path.Join(DIR_SERVER, id, "eula.txt"))
	if err != nil {
		return false
	}
	eulaRegex := regexp.MustCompile(`(?s)eula=true(?s)`)

	return eulaRegex.Match(data)
}

func (gasr *GetAllServerResponse) isPortAlreadyUsed(id string) bool {
	data, err := ioutil.ReadFile(path.Join(DIR_SERVER, id, "msa.std"))
	if err != nil {
		return false
	}

	reg := regexp.MustCompile(`(?s)\[Server thread\/WARN\]: \*\*\*\* FAILED TO BIND TO PORT!(?s)`)

	return reg.MatchString(string(data))
}

func (gasr *GetAllServerResponse) IsRunningOnlineMode(id string) bool {
	data, err := ioutil.ReadFile(path.Join(DIR_SERVER, id, "msa.std"))
	if err != nil {
		return false
	}

	reg := regexp.MustCompile(`(?s)\[Server thread\/WARN\]: \*\*\*\* SERVER IS RUNNING IN OFFLINE\/INSECURE MODE!(?s)`)

	return !reg.MatchString(string(data))
}

func (gasr *GetAllServerResponse) GetUsedWorldName(id string) string {
	data, err := ioutil.ReadFile(path.Join(DIR_SERVER, id, "msa.std"))
	if err != nil {
		return ""
	}

	reg := regexp.MustCompile(`(?s)\[Server thread\/INFO\]: Preparing level "(.*?)"(?s)`)

	return reg.FindString(string(data))
}
