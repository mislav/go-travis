package commands

import (
	"bytes"
	"os"
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
	os.Setenv("TRAVIS_ENDPOINT", configuration.GetDefaultTravisEndpoint())
	os.Setenv("TRAVIS_TOKEN", configuration.GetTravisTokenForEndpoint(endpoint))
	whoamiCmd(&cmd)

	t.Error("Output: " + outBuffer.String())
	t.Error("Error:  " + errBuffer.String())
}
