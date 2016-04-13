package commands

import (
	"github.com/HPI-BP2015H/go-travis-1/client"
	"github.com/HPI-BP2015H/go-travis-1/config"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("builds", buildsCmd)
}

type Branches struct {
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Name       string      `json:"name"`
	LastBuild  Build       `json:"build"`
	Repository *Repository `json:"repo"`
}

type Repository struct {
	Name string `json:"name"`
}

func branchesCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": config.RepoSlug(),
		//"build.event_type": "push",
		//"limit":            "10",
	}

	res, err := client.Travis().PerformAction("branches", "find", params)
	if err != nil {
		panic(err)
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	branches := Branches{}
	res.Unmarshal(&branches)

	for _, branch := range branches.Branches {
		cmd.Stdout.Printf("%s: #%s %s\n", branch.Name, branch.LastBuild.Number, branch.LastBuild.State)
	}
}
