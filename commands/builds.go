package commands

import (
	"strings"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "builds",
			Info:     "alias for history",
			Function: buildsCmd,
		},
		cli.Command{
			Name:     "history",
			Info:     "displays a projects build history",
			Function: buildsCmd,
		},
	)
}

func buildsCmd(cmd *cli.Cmd) cli.ExitValue {
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug":  env.Repo,
		"build.event_type": "push",
		"limit":            "10",
	}

	res, err := env.Client.PerformAction("builds", "find", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		return cli.Failure
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		return cli.Failure
	}

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		printBuild(build, cmd)
	}
	return cli.Success
}

func printBuild(build Build, cmd *cli.Cmd) {
	commitMessage := strings.Replace(build.Commit.Message, "\n", " ", -1)
	if build.HasPassed() {
		cmd.Stdout.PushColor("boldgreen")
	} else {
		cmd.Stdout.PushColor("boldred")
	}
	cmd.Stdout.Print("#" + build.Number + " " + build.State)
	cmd.Stdout.PopColor()
	cmd.Stdout.Cprintf(" %C(yellow)(%s) %C(reset)", build.Branch.Name)
	cmd.Stdout.Println(commitMessage)
}
