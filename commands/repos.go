package commands

import (
	"github.com/fatih/color"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("repos", reposCmd)
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
