package commands

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/run"
)

func TestWhoAmI(t *testing.T) {

	// set up pipes for stderr and stdout
	errPipeRead, errPipeWrite, errErr := os.Pipe()
	if errErr != nil {
		t.Error("Could not create pipe for capturing stderr.")
	}
	outPipeRead, outPipeWrite, outErr := os.Pipe()
	if outErr != nil {
		t.Error("Could not create pipe for capturing stdout.")
	}
	tmpErr := os.Stderr
	tmpOut := os.Stdout
	defer func() {
		os.Stderr = tmpErr
		os.Stdout = tmpOut
	}()
	os.Stderr = errPipeWrite
	os.Stdout = outPipeWrite

	// CLI call
	os.Args = []string{"go-travis", "whoami", "--token", "wrongtoken", "--api-endpoint", "https://api.travis-ci.org"}
	traviscli.Run(client.Travis)

	// close pipes
	errPipeWrite.Close()
	outPipeWrite.Close()

	// read from pipes
	stderrBytes, errErr2 := ioutil.ReadAll(errPipeRead)
	if errErr2 != nil {
		t.Error("Could not read from pipe for capturing stderr.")
	}
	stdoutBytes, outErr2 := ioutil.ReadAll(outPipeRead)
	if outErr2 != nil {
		t.Error("Could not read from pipe for capturing stdout.")
	}
	stderr := string(stderrBytes)
	stdout := string(stdoutBytes)

	// assertions
	if !strings.Contains(stderr, "You need to be logged in to do this.") {
		t.Error("Incorrect error message: " + stderr)
	}
	if stdout != "" {
		t.Error("Output not empty: " + stdout)
	}

}
