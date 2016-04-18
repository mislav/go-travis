package commands

import (
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("whatsup", "lists most recent builds", whatsupCmd)
}

func whatsupCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"include": "branch.last_build",
	}
	repositories, err := GetAllRepositories(params)
	if err != nil {
		color.Red("Error: Could not get Repositories.")
		cmd.Exit(1)
	}
	for _, repo := range repositories.Repositories {
		if repo.Active && (repo.DefaultBranch.LastBuild != nil) {
			printRepoStatus(repo)
		}
	}
}

func printRepoStatus(repo Repository) {
	c := color.New(color.FgRed).PrintfFunc()
	cb := color.New(color.FgRed, color.Bold).PrintfFunc()
	build := repo.DefaultBranch.LastBuild
	if build.State == "passed" {
		c = color.New(color.FgGreen).PrintfFunc()
		cb = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}
	cb("%s ", repo.Slug)
	c(build.State+": #%s \n", build.Number)
}
