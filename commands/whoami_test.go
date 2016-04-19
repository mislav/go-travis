package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func TestWhoAmI(t *testing.T) {
	var outBuffer, errBuffer bytes.Buffer
	cmd := cli.Cmd{
		Stdout: cli.NewWriter(&outBuffer),
		Stderr: cli.NewWriter(&errBuffer),
	}
	configuration := config.DefaultConfiguration()
	endpoint := configuration.GetDefaultTravisEndpoint()
	os.Setenv("TRAVIS_ENDPOINT", endpoint) // TODO: Set in cmd Env
	//os.Setenv("TRAVIS_TOKEN", configuration.GetTravisTokenForEndpoint(endpoint) )
	os.Setenv("TRAVIS_TOKEN", "wrongToken") // TODO: Set in cmd Env
	whoamiCmd(&cmd)

	if !strings.Contains(errBuffer.String(), "Error") {
		t.Error("Output: " + outBuffer.String())
		t.Error("Error:  " + errBuffer.String())
	}

}
