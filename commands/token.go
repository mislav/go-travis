package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
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

func tokenCmd(cmd *cli.Cmd) int {
	if NotLoggedIn(cmd) {
		return 1
	}
	env := cmd.Env.(config.TravisCommandConfig)
	if len(env.Token) > 0 {
		cmd.Stdout.Print("Your access token for ")
		cmd.Stdout.Cprint("yellow", env.Endpoint)
		cmd.Stdout.Print(" is ")
		cmd.Stdout.Cprintln("boldgreen", env.Token)
		return 0
	} else {
		// cmd.Stderr.Println("Not logged in for " + env.Endpoint + ", please run travis login.")
		return 1
	}
}
