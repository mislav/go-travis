package commands

import (
	"fmt"

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
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   string  `json:"description"`
	Active        bool    `json:"active"`
	Private       bool    `json:"private"`
	Owner         *Owner  `json:"owner"`
	DefaultBranch *Branch `json:"default_branch"`
}

func (r *Repository) HasDescription() bool {
	return r.Description != ""
}

type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"login"`
}

func reposCmd(cmd *cli.Cmd) {
	repositories, err := GetAllRepositories()
	if err != nil {
		color.Red("Error: Could not get Repositories.")
		cmd.Exit(1)
	}
	for _, repo := range repositories.Repositories {
		printRepo(repo)
	}
}

//GetAllRepositories returns all the repositories (also those not active in travis)
//of the currently logged in user
func GetAllRepositories() (Repositories, error) {
	params := map[string]string{}
	repositories := Repositories{}
	res, err := client.Travis().PerformAction("repositories", "for_current_user", params)
	defer res.Body.Close()
	if err != nil {
		return repositories, err
	}
	if res.StatusCode > 299 {
		return repositories, fmt.Errorf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
	}
	res.Unmarshal(&repositories)
	return repositories, nil
}

func GetAllRepositoriesWithLastBuild() (Repositories, error) {
	params := map[string]string{
		"include": "branch.last_build",
	}
	repositories := Repositories{}
	res, err := client.Travis().PerformAction("repositories", "for_current_user", params)
	defer res.Body.Close()
	if err != nil {
		return repositories, err
	}
	if res.StatusCode > 299 {
		return repositories, fmt.Errorf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
	}
	res.Unmarshal(&repositories)
	return repositories, nil
}

func printRepo(repo Repository) {
	y := color.New(color.FgYellow, color.Bold).PrintfFunc()
	y(repo.Slug)
	color.Yellow(" (active: %v, private: %v)", repo.Active, repo.Private)
	if repo.HasDescription() {
		color.Green("   Description: %s ", repo.Description)
	}
	println("")
}
