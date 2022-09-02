package server

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/Bearaujus/minecraft-server-api/internal/model"
	"github.com/Bearaujus/minecraft-server-api/pkg"

	"github.com/google/uuid"
)

func (sr *serverResource) GetAllServerResource() (map[string]*model.Server, error) {
	return sr.serverdata, nil
}

func (sr *serverResource) CreateServerResource() (string, error) {
	var resID = uuid.New().String()
	if err := os.MkdirAll(path.Join(model.DIR_SERVER, resID), os.ModePerm); err != nil {
		return "", err
	}

	if err := sr.addServer(resID); err != nil {
		return "", err
	}

	return resID, nil
}

func (sr *serverResource) DeleteServerResource(id string) error {
	srv, err := sr.getServer(id)
	if err != nil {
		return err
	}

	if srv != nil {
		return errors.New("server is running")
	}

	if err := pkg.DeleteDir(path.Join(model.DIR_SERVER, id)); err != nil {
		return err
	}

	delete(sr.serverdata, id)

	return nil
}

func (sr *serverResource) AgreeEulaServerResource(id string) error {
	_, err := sr.getServer(id)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path.Join(model.DIR_SERVER, id, "eula.txt"))
	if err != nil {
		eulaModel := []string{
			`#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://aka.ms/MinecraftEULA).`,
			fmt.Sprintf("#%v", time.Now().Format("Mon Jan 03:04:05 MST 2006")),
			"eula=true",
			"",
		}

		if err := ioutil.WriteFile(path.Join(model.DIR_SERVER, id, "eula.txt"), []byte(strings.Join(eulaModel, "\n")), 0644); err != nil {
			return err
		}

		return nil
	}

	eulaRegex := regexp.MustCompile(`(?s)eula=false(?s)`)
	if !eulaRegex.Match(data) {
		return errors.New("eula already agreed")
	}

	if err := ioutil.WriteFile(path.Join(model.DIR_SERVER, id, "eula.txt"), eulaRegex.ReplaceAll(data, []byte("eula=true")), 0644); err != nil {
		return err
	}

	return nil
}

func (sr *serverResource) StartServerResource(id string, ramGB, port int, worldName string) error {
	srv, err := sr.getServer(id)
	if err != nil {
		return err
	}

	if srv != nil {
		return errors.New("server already started")
	}

	if err := pkg.DeleteDir(path.Join(model.DIR_SERVER, id, "msa.std")); err != nil {
		return err
	}

	fileOut, err := os.OpenFile(path.Join(model.DIR_SERVER, id, "msa.std"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	bufStdout := bufio.NewWriter(fileOut)

	cmdArgs := []string{
		"-XX:+UseG1GC",
		"-XX:+ParallelRefProcEnabled",
		"-XX:MaxGCPauseMillis=200",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:+DisableExplicitGC",
		"-XX:+AlwaysPreTouch",
		"-XX:G1NewSizePercent=30",
		"-XX:G1MaxNewSizePercent=40",
		"-XX:G1HeapRegionSize=8M",
		"-XX:G1ReservePercent=20",
		"-XX:G1HeapWastePercent=5",
		"-XX:G1MixedGCCountTarget=4",
		"-XX:InitiatingHeapOccupancyPercent=15",
		"-XX:G1MixedGCLiveThresholdPercent=90",
		"-XX:G1RSetUpdatingPauseTimePercent=5",
		"-XX:SurvivorRatio=32",
		"-XX:+PerfDisableSharedMem",
		"-XX:MaxTenuringThreshold=1",
		"-Dusing.aikars.flags=https://mcflags.emc.gs",
		"-Daikars.new.flags=true",

		// set ram useage
		fmt.Sprintf("-Xms%vG", ramGB),
		fmt.Sprintf("-Xmx%vG", ramGB),

		// set jar file
		"-jar", "../../jar/server-1.19.2.jar",

		// set port
		"--port", fmt.Sprint(port),

		"--nogui",
	}

	if worldName != "" {
		cmdArgs = append(cmdArgs, "--world", fmt.Sprint(worldName))
	}

	cmd := exec.Command("java", cmdArgs...)
	cmd.Dir = path.Join(model.DIR_SERVER, id)
	cmd.Stdout = bufStdout
	cmd.Stderr = bufStdout
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	modelServer := model.Server{
		ID:        id,
		Port:      port,
		RamGB:     ramGB,
		Cmd:       cmd,
		StdinPipe: &stdinPipe,
		FileOut:   fileOut,
	}

	if err := modelServer.Cmd.Start(); err != nil {
		return err
	}

	if err := sr.setServer(id, &modelServer); err != nil {
		return err
	}

	// if process wasn't started properly, kill it
	go func() error {
		sr.serverdata[id].IsAttemptedToStart = true
		waitTime := time.Second * 120
		tickerTime := time.Millisecond * 500
		ticker := time.NewTicker(tickerTime)

		for range ticker.C {
			srv, err := sr.getServer(id)
			if err != nil {
				return err
			}

			if srv == nil {
				return errors.New("server is not started")
			}

			if _, err := fmt.Fprint(*srv.StdinPipe, ""); err != nil {
				srv.Cmd.Process.Kill()
				sr.serverdata[id] = nil
				return errors.New("server is not started")
			}

			// stop go routine from unexpected error
			waitTime = waitTime - tickerTime
			if waitTime <= 0 {
				srv.Cmd.Process.Kill()
				sr.serverdata[id] = nil
				return errors.New("server took too much time when starting")
			}

			// read data
			data, err := ioutil.ReadFile(path.Join(model.DIR_SERVER, id, "msa.std"))
			if err != nil {
				return err
			}

			// success regex
			regSuccess := regexp.MustCompile(`(?s)\[Server thread\/INFO\]: Done \((.*?)\)! For help, type "help"(?s)`)
			if regSuccess.MatchString(string(data)) {
				sr.serverdata[id].IsAttemptedToStart = false
				break
			}

			// fail regex
			regFailToBindPort := regexp.MustCompile(`(?s)\[Server thread\/WARN\]: \*\*\*\* FAILED TO BIND TO PORT!(?s)`)
			if regFailToBindPort.MatchString(string(data)) {
				srv.Cmd.Process.Kill()
				sr.serverdata[id] = nil
				return errors.New("fail to bind port")
			}
		}

		return nil
	}()

	return nil
}

func (sr *serverResource) StopServerResource(id string) error {
	srv, err := sr.getServer(id)
	if err != nil {
		return err
	}

	if srv == nil {
		return errors.New("server is not started")
	}

	if srv.IsAttemptedToStop {
		return errors.New("server already attempted to stop")
	}

	if _, err := fmt.Fprintln(*srv.StdinPipe, "stop"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(srv.FileOut, "stop"); err != nil {
		return err
	}

	// wait for pipe until broken
	go func() {
		sr.serverdata[id].IsAttemptedToStop = true
		tickerTime := time.Millisecond * 500
		ticker := time.NewTicker(tickerTime)
		for range ticker.C {
			if _, err := fmt.Fprint(*srv.StdinPipe, ""); err != nil {
				srv.Cmd.Process.Kill()
				sr.serverdata[id] = nil
				break
			}
		}
	}()

	return nil
}

func (sr *serverResource) GetServerConsoleResource(id string) ([]byte, error) {
	srv, err := sr.getServer(id)
	if err != nil {
		return nil, err
	}

	if srv == nil {
		return nil, errors.New("server is not started")
	}

	fileOut, err := os.OpenFile(path.Join(model.DIR_SERVER, id, "msa.std"), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fileOut.Close()

	return ioutil.ReadAll(fileOut)
}

func (sr *serverResource) AddServerConsoleResource(id, command string) error {
	srv, err := sr.getServer(id)
	if err != nil {
		return err
	}

	if srv == nil {
		return errors.New("server is not started")
	}

	if command == "stop" {
		return sr.StopServerResource(id)
	}

	if _, err := fmt.Fprintln(*srv.StdinPipe, command); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(srv.FileOut, command); err != nil {
		return err
	}

	return nil
}
