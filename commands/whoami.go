package commands

import (
	"github.com/HPI-BP2015H/go-travis/commands/helper"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("whoami", "outputs the current user", whoamiCmd)
}

func whoamiCmd(cmd *cli.Cmd) {
	user, err := user.CurrentUser()
	if err != nil {
		color.Red("Error: Could not get the current user! \n" + err.Error())
		return
	}
	printUser(user)
}

func printUser(user user.User) {
	g := color.New(color.FgGreen).PrintfFunc()
	gb := color.New(color.FgGreen, color.Bold).PrintfFunc()
	g("You are ")
	gb(user.Login)
	if (user.Name != user.Login) && (user.Name != "") {
		color.Green(" (%s).", user.Name)
	} else {
		color.Green(".")
	}
}
