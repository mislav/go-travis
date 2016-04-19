package commands

import (
	"io"
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "enable",
			Help:     "enables a project",
			Function: enableCmd,
		},
	)
}

func enableCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": os.Getenv("TRAVIS_REPO"),
	}

	res, err := client.Travis().PerformAction("repository", "enable", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		cmd.Exit(1)
	}
	defer res.Body.Close()
	io.Copy(cmd.Stdout, res.Body)
	cmd.Exit(0)
}
