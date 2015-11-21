package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/phayes/hookserve/hookserve"
)

var (
	pushPort          int
	pushNotifHTTPPath string
	pushGithubSecret  string
	projectDir        string
)

func main() {
	var err error
	if pushPort, err = strconv.Atoi(os.Getenv("PUSH_PORT")); err != nil {
		log.Fatal("Env PUSH_PORT must be integer")
	}

	pushNotifHTTPPath = os.Getenv("PUSH_HTTP_PATH")
	if pushNotifHTTPPath == "" {
		log.Fatal("Env PUSH_HTTP_PATH must be present")
	}

	pushGithubSecret = os.Getenv("PUSH_GITHUB_SECRET")

	projectDir = os.Getenv("PUSH_PROJECT_DIR")
	if projectDir == "" {
		log.Fatal("Env PUSH_PROJECT_DIR must be present")
	}

	server := hookserve.NewServer()
	server.Port = pushPort
	server.Path = pushNotifHTTPPath
	server.Secret = pushGithubSecret
	server.GoListenAndServe()

	cliExecuter := cliRunner{os.Stderr}

	fmt.Printf("Starting server...\nExpecting push notification on :%d%s\nProject path %s\n", pushPort, pushNotifHTTPPath, projectDir)

	for {
		select {
		case event := <-server.Events:
			err := pull(event, &cliExecuter)
			if err != nil {
				log.Fatal(err)
			}

			err = build(&cliExecuter)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

type runner interface {
	Run(cmd *exec.Cmd) error
}

type cliRunner struct {
	log io.Writer
}

func (cli *cliRunner) Run(cmd *exec.Cmd) error {
	cmd.Stderr = cli.log

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func pull(event hookserve.Event, executer runner) error {
	cliCmd := exec.Command("git", "pull")
	cliCmd.Dir = projectDir

	return executer.Run(cliCmd)
}

func build(executer runner) error {
	cmd := exec.Command("hugo")
	cmd.Dir = projectDir

	return executer.Run(cmd)
}
