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

func statusCmd(cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug": env.Repo,
		"limit":           "1",
		"sort_by":         "id:desc",
	}

	res, err := env.Client.PerformAction("builds", "find", params)
	if err != nil {
		cmd.Stderr.Println("Build not found.")
		return
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) > 0 {
		printStatus(builds.Builds[0], cmd)
	} else {
		cmd.Stderr.Println("This repository has no builds yet.")
	}
	cmd.Exit(0)
}

func printStatus(build Build, cmd *cli.Cmd) {
	cmd.Stdout.Print("Build #" + build.Number)
	PushColorAccordingToBuildStatusBold(build, cmd)
	cmd.Stdout.Printf(" %s \n", build.State)
	cmd.Stdout.PopColor()
}
