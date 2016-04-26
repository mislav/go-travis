package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "status",
			Info:     "checks status of the latest build",
			Function: statusCmd,
		},
	)
}

func statusCmd(cmd *cli.Cmd) cli.ExitValue {
	if NoRepoSpecified(cmd) {
		return cli.Failure
	}
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug": env.Repo,
		"limit":           "1",
		"sort_by":         "id:desc",
	}

	res, err := env.Client.PerformAction("builds", "find", params)
	if err != nil {
		cmd.Stderr.Println("Build not found.")
		return cli.Failure
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		return cli.Failure
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) == 0 {
		cmd.Stderr.Println("This repository has no builds yet.")
		return cli.Failure
	}
	printStatus(builds.Builds[0], cmd)
	return cli.Success
}

func printStatus(build Build, cmd *cli.Cmd) {
	cmd.Stdout.Print("Build #" + build.Number)
	PushBoldColorAccordingToBuildStatus(build, cmd)
	cmd.Stdout.Printf(" %s \n", build.State)
	cmd.Stdout.PopColor()
}
