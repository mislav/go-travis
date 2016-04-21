package commands

import (
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

func whoamiCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) {
		return cli.Failure
	}
	env := cmd.Env.(config.TravisCommandConfig)
	user, _ := CurrentUser(env.Client)
	printUser(user, cmd)
	return cli.Success
}

func printUser(user User, cmd *cli.Cmd) {
	cmd.Stdout.PushColor("green")
	cmd.Stdout.Print("You are ")
	cmd.Stdout.Cprint("boldgreen", user.Login)
	if (user.Name != user.Login) && (user.Name != "") {
		cmd.Stdout.Printf(" (%s)", user.Name)
	}
	cmd.Stdout.Print(".\n")
	cmd.Stdout.PopColor()
}
