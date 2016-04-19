package commands

import (
	"os"
	"strings"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "builds",
			Help:     "alias for history",
			Function: buildsCmd,
		},
	)
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "history",
			Help:     "displays a projects build history",
			Function: buildsCmd,
		},
	)
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number     string  `json:"number"`
	State      string  `json:"state"`
	StartedAt  string  `json:"started_at"`
	FinishedAt string  `json:"finished_at"`
	Duration   int     `json:"duration"`
	EventType  string  `json:"event_type"`
	Branch     *Branch `json:"branch"`
	Commit     *Commit `json:"commit"`
	Jobs       Jobs    `json:"jobs"`
}

type Commit struct {
	Message string `json:"message"`
}

type Jobs struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Number string `json:"number"`
	State  string `json:"state"`
}

func buildsCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug":  os.Getenv("TRAVIS_REPO"),
		"build.event_type": "push",
		"limit":            "10",
	}

	res, err := client.Travis().PerformAction("builds", "find", params)
	if err != nil {
		panic(err)
	}
	if res.StatusCode > 299 {
		color.Red("Unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		printBuild(build)
	}
}

func printBuild(build Build) {
	commitMessage := strings.Replace(build.Commit.Message, "\n", " ", -1)
	y := color.New(color.FgYellow).PrintfFunc()
	c := color.New(color.FgRed, color.Bold).PrintfFunc()
	if build.State == "passed" {
		c = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}
	c("#%s %s ", build.Number, build.State)
	y("(%s) ", build.Branch.Name)
	print(commitMessage + "\n")
}
