package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cmd := cli.Command{
		Name:     "endpoint",
		Help:     "displays or changes the API endpoint",
		Function: endpointCmd,
	}
	cmd.RegisterFlag(
		cli.Flag{
			Long:  "--set-default",
			Ftype: false,
		},
	)
	cmd.RegisterFlag(
		cli.Flag{
			Long:  "--drop-default",
			Ftype: false,
		},
	)
	cli.AppInstance().RegisterCommand(cmd)
}

func endpointCmd(cmd *cli.Cmd) {

	configuration := config.DefaultConfiguration()
	endpoint := cmd.Env["TRAVIS_ENDPOINT"]

	if cmd.Flags.IsProvided("--set-default") {
		configuration.StoreDefaultTravisEndpoint(endpoint)
		cmd.Stdout.Cprintln("green", "Stored default Travis endpoint")
	}

	if cmd.Flags.IsProvided("--drop-default") {
		configuration.DeleteDefaultTravisEndpoint()
		cmd.Stdout.Cprintln("green", "Deleted default Travis endpoint")
	}

	println("API endpoint: " + endpoint)
	cmd.Exit(0)
}
