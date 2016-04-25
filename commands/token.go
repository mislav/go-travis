package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "token",
			Info:     "outputs the secret API token",
			Function: tokenCmd,
		},
	)
}

func tokenCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) {
		return cli.Failure
	}
	env := cmd.Env.(config.TravisCommandConfig)
	if len(env.Token) > 0 {
		cmd.Stdout.Print("Your access token for ")
		cmd.Stdout.Cprint("yellow", env.Endpoint)
		cmd.Stdout.Print(" is ")
		cmd.Stdout.Cprintln("boldgreen", env.Token)
		return cli.Success
	} else {
		// cmd.Stderr.Println("Not logged in for " + env.Endpoint + ", please run travis login.")
		return cli.Failure
	}
}
