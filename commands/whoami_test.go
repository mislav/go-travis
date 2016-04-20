package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/HPI-BP2015H/go-travis/config"
)

func TestWhoAmI(t *testing.T) {

	//create cmd and redirect stdout and stderr
	cmd := exec.Command("go-travis", "whoami")
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	//set cmd's env
	env := os.Environ()
	endpoint := config.DefaultConfiguration().GetDefaultTravisEndpoint()
	env = append(env, fmt.Sprintf("TRAVIS_ENDPOINT=%s", endpoint))
	env = append(env, fmt.Sprintf("TRAVIS_TOKEN=%s", "wrongToken"))
	cmd.Env = env

	cmd.Run()

	if !strings.Contains(err.String(), "Error") {
		t.Error("Output: " + out.String())
		t.Error("Error:  " + err.String())
	}

}
