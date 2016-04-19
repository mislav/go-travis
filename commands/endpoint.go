package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("endpoint", "displays or changes the API endpoint", endpointCmd)
}

func endpointCmd(cmd *cli.Cmd) {
	setDefaultFlag, args := cmd.Args.ExtractFlag("", "--set-default", false)
	dropDefaultFlag, args := args.ExtractFlag("", "--drop-default", false)

	configuration := config.DefaultConfiguration()
	endpoint := os.Getenv("TRAVIS_ENDPOINT")

	if setDefaultFlag.IsProvided() {
		configuration.StoreDefaultTravisEndpoint(endpoint)
		cmd.Stdout.Cprintln("green", "Stored default Travis endpoint")
	}

	if dropDefaultFlag.IsProvided() {
		configuration.DeleteDefaultTravisEndpoint()
		cmd.Stdout.Cprintln("green", "Deleted default Travis endpoint")
	}

	println("API endpoint: " + endpoint)
	cmd.Exit(0)
}
