package commands

import (
	"github.com/HPI-BP2015H/go-travis/commands/helper"
	"github.com/HPI-BP2015H/go-travis/config"
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
	env := cmd.Env.(config.TravisCommandConfig)

	user, err := user.CurrentUser(env.Client)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get the current user!")
		cmd.Exit(1)
	}
	printUser(user, cmd)
	cmd.Exit(0)
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
