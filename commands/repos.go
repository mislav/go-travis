package commands

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("repos", "lists repositories the user has certain permissions on", reposCmd)
}

type Repositories struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Private     bool   `json:"private"`
	Owner       *Owner `json:"owner"`
}

func (r *Repository) HasDescription() bool {
	return r.Description != ""
}

type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"login"`
}

func reposCmd(cmd *cli.Cmd) {
	params := map[string]string{}
	res, err := client.Travis().PerformAction("repositories", "for_current_user", params)
	defer res.Body.Close()
	if err != nil {
		color.Red("Error: Could not connect to Travis! \n" + err.Error())
		color.Yellow("Fall back to asking Github:")
		reposGithub()
		return
	}
	if res.StatusCode > 299 {
		color.Red("Error: Unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	repositories := Repositories{}
	res.Unmarshal(&repositories)
	for _, repo := range repositories.Repositories {
		printRepoColorful(repo)
	}

}

func printRepoColorful(repo Repository) {

	y := color.New(color.FgYellow, color.Bold).PrintfFunc()
	y(repo.Slug)
	color.Yellow(" (active: %v, private: %v)", repo.Active, repo.Private)
	if repo.HasDescription() {
		color.Green("Description: %s ", repo.Description)
	}
	println("")

}

func reposGithub() {
	github, _ := LoginToGitHub("")
	repos, _, err := github.Repositories.List("", nil)
	if err != nil {
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	println("These are your current Repositories:")
	for _, repo := range repos {
		color.Blue("     " + *repo.FullName)
	}
}
