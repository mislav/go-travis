package commands

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestHelpCmd(t *testing.T) {
	cmd := exec.Command("go-travis", "help")
	//cmd.Stdin = strings.NewReader("some input") //if we need to supply input to the command
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("THIS IS OUTPUT!!!: %s\n", out.String())

}
