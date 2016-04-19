package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/HPI-BP2015H/go-utils/cli"
)

func TestHelpCmd(t *testing.T) {
	var outBuffer, errBuffer bytes.Buffer
	cmd := cli.Cmd{
		Stdout: cli.NewWriter(&outBuffer),
		Stderr: cli.NewWriter(&errBuffer),
	}

	helpCmd(&cmd)

	if !strings.Contains(outBuffer.String(), "Available commands:") {
		t.Error("Output: " + outBuffer.String())
		t.Error("Error:  " + errBuffer.String())
	}

}
