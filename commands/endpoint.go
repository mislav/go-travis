package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/fatih/color"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("endpoint", endpointCmd)
}

func endpointCmd(cmd *cli.Cmd) {
	setDefaultFlag, args := cmd.Args.ExtractFlag("", "--set-default", false)
	dropDefaultFlag, args := args.ExtractFlag("", "--drop-default", false)

	configuration := client.DefaultConfiguration()
	endpoint := os.Getenv("TRAVIS_ENDPOINT")

	if setDefaultFlag.IsProvided() {
		configuration.StoreDefaultTravisEndpoint(endpoint)
		color.Green("Stored default Travis endpoint" + "\n")
	}

	if dropDefaultFlag.IsProvided() {
		configuration.DeleteDefaultTravisEndpoint()
		color.Green("Deleted default Travis endpoint" + "\n")
	}

	println("API endpoint: " + endpoint)

}
