package commands

import (
	"strings"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/fatih/color"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("builds", buildsCmd)
	cli.Register("history", buildsCmd)
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number string  `json:"number"`
	State  string  `json:"state"`
	Branch *Branch `json:"branch"`
	Commit *Commit `json:"commit"`
}

type Commit struct {
	Message string `json:"message"`
}

func buildsCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug":  config.RepoSlug(),
		"build.event_type": "push",
		"limit":            "10",
	}

	res, err := client.Travis().PerformAction("builds", "find", params)
	if err != nil {
		panic(err)
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		printBuildColorful(build)
	}
}

func printBuildColorful(build Build) {
	commitMessage := strings.Replace(build.Commit.Message, "\n", " ", -1)
	y := color.New(color.FgYellow).PrintfFunc()
	c := color.New(color.FgRed, color.Bold).PrintfFunc()
	if build.State == "passed" {
		c = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}
	c("#%s %s  ", build.Number, build.State)
	y("(%s) ", build.Branch.Name)
	print(commitMessage + "\n")
}
