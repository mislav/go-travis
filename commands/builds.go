package commands

import (
	"os"
	"strings"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("builds", "alias for history", buildsCmd)
	cli.Register("history", "displays a projects build history", buildsCmd)
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

func (b *Build) HasPassed() bool {
	return b.State == "passed"
}

func (b *Build) IsNotYetFinished() bool {
	return ((b.State == "created") || (b.State == "started"))
}

func PushColorAccordingToBuildStatusBold(build Build, cmd *cli.Cmd) {
	if build.HasPassed() {
		cmd.Stdout.PushColor("boldgreen")
	} else {
		if build.IsNotYetFinished() {
			cmd.Stdout.PushColor("boldyellow")
		} else {
			cmd.Stdout.PushColor("boldred")
		}
	}
}

func PushColorAccordingToBuildStatus(build Build, cmd *cli.Cmd) {
	if build.HasPassed() {
		cmd.Stdout.PushColor("green")
	} else {
		if build.IsNotYetFinished() {
			cmd.Stdout.PushColor("yellow")
		} else {
			cmd.Stdout.PushColor("red")
		}
	}
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
		cmd.Stderr.Println(err.Error())
		cmd.Exit(1)
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		printBuild(build, cmd)
	}
	cmd.Exit(0)
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
	cmd.Stdout.Cprint("yellow", "(%s) ", build.Branch.Name)
	cmd.Stdout.Println(commitMessage)
}
