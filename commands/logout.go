package commands

import (
	"github.com/HPI-BP2015H/go-travis/commands/helper"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "logout",
			Help:     "deletes the stored API token",
			Function: logoutCmd,
		},
	)
}

func logoutCmd(cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)
	user, _ := user.CurrentUser(env.Client)
	config := config.DefaultConfiguration()
	config.DeleteTravisTokenForEndpoint(env.Endpoint)
	cmd.Stdout.Cprint("green", "%s is now logged out.", user)
	cmd.Exit(0)
}
