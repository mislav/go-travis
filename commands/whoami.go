package commands

import (
	"github.com/HPI-BP2015H/go-travis/commands/helper"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "whoami",
			Help:     "outputs the current user",
			Function: whoamiCmd,
		},
	)
}

func whoamiCmd(cmd *cli.Cmd) {
	user, err := user.CurrentUser()
	if err != nil {
		cmd.Stderr.Cprintln("red", "Error: Could not get the current user! \n"+err.Error())
		return
	}
	printUser(user, cmd)
}

func printUser(user user.User, cmd *cli.Cmd) {
	cmd.Stdout.PushColor("green")
	cmd.Stdout.Print("You are ")
	cmd.Stdout.Cprint("boldgreen", user.Login)
	if (user.Name != user.Login) && (user.Name != "") {
		cmd.Stdout.Printf(" (%s)", user.Name)
	}
	cmd.Stdout.Print(".\n")
	cmd.Stdout.PopColor()
}
