package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	_ "github.com/mislav/go-travis/commands"
	"github.com/mislav/go-utils/cli"
	"github.com/mislav/go-utils/pathname"
)

func main() {
	args := cli.NewArgs(os.Args)
	cmdName := args.Peek(0)
	if cmdName == "" {
		cmdName = "builds"
	}

	repoFlag, args := args.ExtractFlag("-r", "--repo", "REPOSITORY")
	tokenFlag, args := args.ExtractFlag("-t", "--token", "TOKEN")
	debugFlag, args := args.ExtractFlag("", "--debug", false)

	if repoFlag.IsProvided() {
		os.Setenv("TRAVIS_REPO", repoFlag.String())
	}
	if tokenFlag.IsProvided() {
		os.Setenv("TRAVIS_TOKEN", tokenFlag.String())
	}
	if debugFlag.IsProvided() {
		if debugFlag.Bool() {
			os.Setenv("TRAVIS_DEBUG", "1")
		} else {
			os.Setenv("TRAVIS_DEBUG", "")
		}
	}

	cmdFunc := cli.Lookup(cmdName)
	if cmdFunc != nil {
		cmd := cli.NewCmd(args.SubcommandArgs(cmdName))
		cmdFunc(cmd)
	} else {
		exeName := args.ProgramName() + "-" + cmdName
		results := pathname.FindInPath(exeName, strings.Split(os.Getenv("PATH"), ":"))

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
