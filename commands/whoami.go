package commands

import (
	"github.com/fatih/color"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("whoami", whoamiCmd)
}

func whoamiCmd(cmd *cli.Cmd) {
	github := Login()
	user, _, err := github.Users.Get("")
	if err != nil {
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	color.Green("You are logged into the account " + *(user.Login) + ".\n")
}
