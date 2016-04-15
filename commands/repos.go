package commands

import (
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("repos", "lists repositories the user has certain permissions on", reposCmd)
}

func reposCmd(cmd *cli.Cmd) {
	github := LoginToGitHub()
	repos, _, err := github.Repositories.List("", nil)
	if err != nil {
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	println("These are your current Repositories:")
	for _, repo := range repos {
		color.Blue("     " + *repo.FullName)
	}
}
