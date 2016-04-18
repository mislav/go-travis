package commands

import (
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("whatsup", "lists most recent builds", whatsupCmd)
}

func whatsupCmd(cmd *cli.Cmd) {
	repositories, err := GetAllRepositories()
	if err != nil {
		color.Red("Error: Could not get Repositories.")
		cmd.Exit(1)
	}
	for _, repo := range repositories.Repositories {
		if repo.Active {
			printRepoStatus(repo)
		}
	}
}

func printRepoStatus(repo Repository) {
	c := color.New(color.FgRed).PrintfFunc()
	cb := color.New(color.FgRed, color.Bold).PrintfFunc()
	color.Yellow(repo.DefaultBranch.LastBuild.Number)
	build := repo.DefaultBranch.LastBuild
	color.Yellow(build.State + build.Number + build.Commit.Message)
	if build.State == "passed" {
		c = color.New(color.FgRed).PrintfFunc()
		cb = color.New(color.FgRed, color.Bold).PrintfFunc()
	}
	cb("%s ", repo.Name)
	c(build.State+": %d \n", build.Number)
}
