package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/phayes/hookserve/hookserve"
)

func init() {
	os.Setenv("PUSH_PORT", "8578")
	os.Setenv("PUSH_HTTP_PATH", "/push")
	os.Setenv("PUSH_PROJECT_DIR", ".")
}

type TestRunner struct {
	output *bytes.Buffer
}

func (r TestRunner) Run(cmd *exec.Cmd) error {
	fmt.Fprintf(r.output, "%s", cmd.Args)
	return nil
}

func TestPull(t *testing.T) {
	var out bytes.Buffer
	executer := TestRunner{&out}

	pull(hookserve.Event{}, executer)

	if out.String() != "[git pull]" {
		t.Errorf("Expected %q got %q", "[git pull]", out.String())
	}
}

func TestBuild(t *testing.T) {
	var out bytes.Buffer
	executer := TestRunner{&out}

	build(executer)

	if out.String() != "[hugo]" {
		t.Errorf("Expected %q got %q", "[hugo]", out.String())
	}
}

func TestRun(t *testing.T) {
	var out bytes.Buffer
	executer := cliRunner{&out}

	cliCmd := exec.Command("echo", "go rules")
	var out2 bytes.Buffer
	cliCmd.Stdout = &out2

	executer.Run(cliCmd)

	if out2.String() != "go rules\n" {
		t.Errorf("Expected %q got %q", "go rules", out2.String())
	}
}
