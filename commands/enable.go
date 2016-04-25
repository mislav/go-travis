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
			Info:     "enables a project",
			Function: enableCmd,
		},
	)
}

func enableCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) || NoRepoSpecified(cmd) {
		return cli.Failure
	}
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"repository.slug": env.Repo,
	}

	res, err := env.Client.PerformAction("repository", "enable", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		return cli.Failure
	}
	defer res.Body.Close()
	io.Copy(cmd.Stdout, res.Body)
	return cli.Success
}
