package commands

import (
	"io"

	"github.com/HPI-BP2015H/go-travis/config"
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
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug": env.Repo,
	}

	res, err := env.Client.PerformAction("repository", "enable", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		cmd.Exit(1)
	}
	defer res.Body.Close()
	io.Copy(cmd.Stdout, res.Body)
	cmd.Exit(0)
}
