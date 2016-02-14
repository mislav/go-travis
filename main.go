package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	_ "github.com/mislav/go-travis/commands"
	"github.com/mislav/go-utils/cli"
	"github.com/mislav/go-utils/utils"
)

func main() {
	args := cli.NewArgs(os.Args)
	cmdName := args.At(0)
	if cmdName == "" {
		cmdName = "builds"
	}

	cmd := cli.Lookup(cmdName)
	if cmd != nil {
		cmd(args.SubcommandArgs(cmdName))
	} else {
		exeName := args.ProgramName() + "-" + cmdName
		results := utils.FindInPath(exeName, strings.Split(os.Getenv("PATH"), ":"))

		if len(results) > 0 {
			exeCmd := results[0]

			argv := []string{exeName}
			argv = append(argv, os.Args[2:]...)

			err := syscall.Exec(exeCmd.String(), argv, os.Environ())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", exeName, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", exeName)
			os.Exit(1)
		}
	}
}
