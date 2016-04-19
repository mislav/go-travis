package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "token",
			Help:     "outputs the secret API token",
			Function: tokenCmd,
		},
	)
}

func tokenCmd(cmd *cli.Cmd) {
	token := os.Getenv("TRAVIS_TOKEN")
	endpoint := os.Getenv("TRAVIS_ENDPOINT")

	if len(token) > 0 {
		cmd.Stdout.Print("Your access token for ")
		cmd.Stdout.Cprint("yellow", endpoint)
		cmd.Stdout.Print(" is ")
		cmd.Stdout.Cprintln("boldgreen", os.Getenv("TRAVIS_TOKEN"))
	} else {
		cmd.Stderr.Println("Not logged in for " + endpoint + ", please run travis login.")
	}
}
