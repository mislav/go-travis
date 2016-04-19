package commands

import (
	"fmt"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "repos",
			Help:     "lists repositories the user has certain permissions on",
			Function: reposCmd,
		},
	)
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
	repositories, err := GetAllRepositories(nil)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get Repositories.")
		cmd.Exit(1)
	}
	for _, repo := range repositories.Repositories {
		printRepo(repo, cmd)
	}
	cmd.Exit(0)
}

//GetAllRepositories returns all the repositories (also those not active in travis)
//of the currently logged in user. also takes params
func GetAllRepositories(params map[string]string) (Repositories, error) {
	if params == nil {
		params = map[string]string{}
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

func printRepo(repo Repository, cmd *cli.Cmd) {
	cmd.Stdout.Cprint("boldyellow", repo.Slug)
	cmd.Stdout.Cprintln("yellow", " (active: %v, private: %v)", repo.Active, repo.Private)
	if repo.HasDescription() {
		cmd.Stdout.Cprint("green", "   Description: %s ", repo.Description)
	}
	println("")
}
