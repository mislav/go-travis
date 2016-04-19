package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/HPI-BP2015H/go-travis/client"
	_ "github.com/HPI-BP2015H/go-travis/commands"
	"github.com/HPI-BP2015H/go-travis/config"

	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/HPI-BP2015H/go-utils/pathname"
)

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

	app := cli.AppInstance()
	app.Name = "go-travis"
	app.DefaultCommandName = "help"

	app.RegisterFlag(
		cli.Flag{
			Short: "-r",
			Long:  "--repo",
			Ftype: "REPOSITORY",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-t",
			Long:  "--token",
			Ftype: "TOKEN",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-e",
			Long:  "--api-endpoint",
			Ftype: "URL",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--org",
			Ftype: false,
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--pro",
			Ftype: false,
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--staging",
			Ftype: false,
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--debug",
			Ftype: false,
		},
	)

	app.Before = func(cmd *cli.Cmd) {
		configuration := config.DefaultConfiguration()

		repo := ""
		if cmd.Flags.IsProvided("--repo") {
			repo = cmd.Flags.String("--repo")
		} else {
			repo = config.RepoSlugFromGit()
		}

		endpoint := configuration.GetDefaultTravisEndpoint()
		if cmd.Flags.IsProvided("--org") {
			endpoint = config.TravisOrgEndpoint
		}
		if cmd.Flags.IsProvided("--pro") {
			endpoint = config.TravisProEndpoint
		}
		if cmd.Flags.IsProvided("--staging") {
			endpoint = config.TravisStagingEndpoint
		}
		if cmd.Flags.IsProvided("--api-endpoint") {
			endpoint = cmd.Flags.String("--api-endpoint")
		}

		debug := cmd.Flags.IsProvided("--debug")

		token := configuration.GetTravisTokenForEndpoint(endpoint)
		if cmd.Flags.IsProvided("--token") {
			token = cmd.Flags.String("--token")
		}

		commandConfig := config.TravisCommandConfig{
			Config:   configuration,
			Repo:     repo,
			Endpoint: endpoint,
			Token:    token,
			Client:   client.Travis(endpoint, token, debug),
			Debug:    debug,
		}
		cmd.Env = commandConfig

	}

	app.Fallback = func(cmd *cli.Cmd, cmdName string) {
		env := cmd.Env.(config.TravisCommandConfig)

		os.Setenv("TRAVIS_REPO", env.Repo)
		os.Setenv("TRAVIS_TOKEN", env.Token)
		os.Setenv("TRAVIS_ENDPOINT", env.Endpoint)
		if env.Debug {
			os.Setenv("TRAVIS_DEBUG", "true")
		}

		exeName := cmd.Args.ProgramName() + "-" + cmdName
		results := pathname.FindInPath(exeName, strings.Split(os.Getenv("PATH"), ":"))

		if len(results) > 0 {
			exeCmd := results[0]

			argv := []string{exeName}
			argv = append(argv, cmd.Args.Slice(1)...)

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

	app.Run(os.Args)

}
