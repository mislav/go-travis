package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	_ "github.com/HPI-BP2015H/go-travis/commands"
	"github.com/HPI-BP2015H/go-travis/config"

	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/HPI-BP2015H/go-utils/pathname"
)

const defaultCommand = "help"

// main the current implementation is not respection the debug flag
// The following arguments from the original travis cli are missing:
// -i, --[no-]interactive           be interactive and colorful
// -E, --[no-]explode               don't rescue exceptions
//     --skip-version-check         don't check if travis client is up to date
//     --skip-completion-check      don't check if auto-completion is set up
// -I, --[no-]insecure              do not verify SSL certificate of API endpoint
//     --debug-http                 show HTTP(S) exchange
// -X, --enterprise [NAME]          use enterprise setup (optionally takes name for multiple setups)
func main() {
	args := cli.NewArgs(os.Args)
	cmdName := args.Peek(0)
	if cmdName == "" {
		cmdName = defaultCommand
	}

	configuration := config.DefaultConfiguration()

	repoFlag, args := args.ExtractFlag("-r", "--repo", "REPOSITORY")
	tokenFlag, args := args.ExtractFlag("-t", "--token", "TOKEN")
	endpointFlag, args := args.ExtractFlag("-e", "--api-endpoint", "URL")
	orgEndpointFlag, args := args.ExtractFlag("", "--org", false)
	proEndpointFlag, args := args.ExtractFlag("", "--pro", false)
	stagingEndpointFlag, args := args.ExtractFlag("", "--staging", false)
	debugFlag, args := args.ExtractFlag("", "--debug", false)

	if repoFlag.IsProvided() {
		os.Setenv("TRAVIS_REPO", repoFlag.String())
	} else {
		os.Setenv("TRAVIS_REPO", config.RepoSlugFromGit())
	}

	endpoint := configuration.GetDefaultTravisEndpoint()
	if orgEndpointFlag.IsProvided() {
		endpoint = config.TravisOrgEndpoint
	}
	if proEndpointFlag.IsProvided() {
		endpoint = config.TravisProEndpoint
	}
	if stagingEndpointFlag.IsProvided() {
		endpoint = config.TravisStagingEndpoint
	}
	if endpointFlag.IsProvided() {
		endpoint = endpointFlag.String()
	}
	os.Setenv("TRAVIS_ENDPOINT", endpoint)

	token := configuration.GetTravisTokenForEndpoint(endpoint)
	if tokenFlag.IsProvided() {
		token = tokenFlag.String()
	}
	os.Setenv("TRAVIS_TOKEN", token)

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
			argv = append(argv, args.Slice(1)...)

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
