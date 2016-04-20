package commands

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestWhoAmI(t *testing.T) {

	//create cmd and redirect stdout and stderr
	cmd := exec.Command("go-travis", "whoami", "--token", "wrongtoken", "--api-endpoint", "https://api.travis-ci.org")
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	cmd.Run()

	if !strings.Contains(err.String(), "You need to be logged in to do this.") {
		t.Error("Output: " + out.String())
		t.Error("Error: " + err.String())
	}

}
