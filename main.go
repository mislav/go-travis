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
			Ftype: "REPOSITORY_SLUG",
			Help:  "the repository on GitHub",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-t",
			Long:  "--token",
			Ftype: "ACCESS_TOKEN",
			Help:  "access token to use",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-e",
			Long:  "--api-endpoint",
			Ftype: "URL",
			Help:  "Travis API server to talk to",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--org",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api.travis-ci.com/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--pro",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api.travis-ci.com/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--staging",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api-staging.travis-ci.org/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--debug",
			Ftype: false,
			Help:  "show API requests",
		},
	)

	app.Before = func(cmd *cli.Cmd) {
		configuration := config.DefaultConfiguration()

		endpoint := configuration.GetDefaultTravisEndpoint()
		if cmd.Parameters.IsProvided("--org") {
			endpoint = config.TravisOrgEndpoint
		}
		if cmd.Parameters.IsProvided("--pro") {
			endpoint = config.TravisProEndpoint
		}
		if cmd.Parameters.IsProvided("--staging") {
			endpoint = config.TravisStagingEndpoint
		}
		endpoint = cmd.Parameters.String("--api-endpoint", endpoint)

		debug := cmd.Parameters.IsProvided("--debug")
		token := cmd.Parameters.String("--token", configuration.GetTravisTokenForEndpoint(endpoint))

		commandConfig := config.TravisCommandConfig{
			Config:   configuration,
			Repo:     cmd.Parameters.String("--repo", config.RepoSlugFromGit()),
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
