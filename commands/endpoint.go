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
			Help:  "store endpoint as global default",
		},
	)
	cmd.RegisterFlag(
		cli.Flag{
			Long:  "--drop-default",
			Ftype: false,
			Help:  "delete stored default endpoint",
		},
	)
	cli.AppInstance().RegisterCommand(cmd)
}

func endpointCmd(cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)

	if cmd.Parameters.IsProvided("--set-default") {
		env.Config.StoreDefaultTravisEndpoint(env.Endpoint)
		cmd.Stdout.Cprintln("green", "Stored default Travis endpoint")
	}

	if cmd.Parameters.IsProvided("--drop-default") {
		env.Config.DeleteDefaultTravisEndpoint()
		cmd.Stdout.Cprintln("green", "Deleted default Travis endpoint")
	}

	println("API endpoint: " + env.Endpoint)
	cmd.Exit(0)
}
