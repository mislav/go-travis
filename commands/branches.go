package commands

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/fatih/color"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("branches", branchesCmd)
}

type Branches struct {
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Name       string      `json:"name"`
	LastBuild  Build       `json:"last_build"`
	Repository *Repository `json:"repo"`
}

type Repository struct {
	Name string `json:"name"`
}

func branchesCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": config.RepoSlug(),
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
		printBranchColorful(branch)
	}
}

func printBranchColorful(branch Branch) {
	color.Yellow("%s:  ", branch.Name)
	c := color.New(color.FgRed, color.Bold).PrintfFunc()
	if branch.LastBuild.State == "passed" {
		c = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}
	c("  #%s  %s\n", branch.LastBuild.Number, branch.LastBuild.State)
}
