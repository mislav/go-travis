package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "status",
			Help:     "checks status of the latest build",
			Function: statusCmd,
		},
	)
}

func statusCmd(cmd *cli.Cmd) int {
	if NoRepoSpecified(cmd) {
		return 1
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
		return 1
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		return 1
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) == 0 {
		cmd.Stderr.Println("This repository has no builds yet.")
		return 1
	}
	printStatus(builds.Builds[0], cmd)
	return 0
}

func printStatus(build Build, cmd *cli.Cmd) {
	cmd.Stdout.Print("Build #" + build.Number)
	PushColorAccordingToBuildStatusBold(build, cmd)
	cmd.Stdout.Printf(" %s \n", build.State)
	cmd.Stdout.PopColor()
}
