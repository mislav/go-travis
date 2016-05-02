package commands

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestHelpCmd(t *testing.T) {
	cmdName := "go-travis"
	cmd := exec.Command(cmdName, "help")
	//cmd.Stdin = strings.NewReader("some input") //if we need to supply input to the command
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	exitError := cmd.Run()
	if exitError != nil {
		t.Error("Help exited not with 0. \nError:" + exitError.Error())
	}

	if !strings.Contains(out.String(), "Run "+cmdName+" help COMMAND for more infos.") {
		t.Error(out.String())
		t.Error(err.String())
	}

}
