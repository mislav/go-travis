package commands

import (
	"io"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("whoami", "outputs the current user", whoamiCmd)
}

func whoamiCmd(cmd *cli.Cmd) {
	params := map[string]string{}
	res, err := client.Travis().PerformAction("user", "current", params)
	defer res.Body.Close()
	if err != nil {
		color.Red("Error: Could not connect to Travis! \n" + err.Error())
		color.Yellow("Fall back to asking Github:")
		whoamiGithub()
		return
	}
	io.Copy(cmd.Stdout, res.Body)
}

func whoamiGithub() {
	github := LoginToGithub("")
	user, _, err := github.Users.Get("")
	if err != nil {
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	color.Green("You are logged into the account " + *(user.Login) + ".\n")
}
