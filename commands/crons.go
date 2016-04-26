package commands

import (
	"sort"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "crons",
			Info:     "Shows all cron jobs of all repositories (may produce large output)",
			Function: cronsCmd,
		},
	)
}

func cronsCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) {
		return cli.Failure
	}

	//get all repos for this user
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"include": "branch.last_build",
	}
	repositories, err := GetAllRepositories(params, env.Client)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get Repositories.")
		return cli.Failure
	}

	//asyncronously fetch all cons for every repo
	var cRes chan *client.Response = make(chan *client.Response)
	var cErr chan error = make(chan error)
	for _, repo := range repositories.Repositories {
		go getCronsForRepo(repo, cRes, cErr, cmd)
	}

	//handle responses
	var responses cronsArray
	for range repositories.Repositories {
		res := <-cRes
		err := <-cErr
		if err != nil {
			cmd.Stderr.Println("Error: Could not get crons!")
			cmd.Stderr.Println(err.Error())
			return cli.Failure
		}
		crons := Crons{}
		res.Unmarshal(&crons)
		if len(crons.Crons) != 0 {
			responses = append(responses, crons)
		}
	}
	if len(responses) == 0 {
		cmd.Stdout.Cprintln("yellow", "None of your repositories has crons!")
		return cli.Success
	}

	//print each repo with crons
	sort.Sort(responses)
	for _, crons := range responses {
		repo := crons.Crons[0].Repository
		if repo.DefaultBranch.LastBuild != nil {
			PushBoldColorAccordingToBuildStatus(*repo.DefaultBranch.LastBuild, cmd)
			cmd.Stdout.Println(repo.Slug)
			cmd.Stdout.PopColor()
		} else {
			cmd.Stdout.Println(repo.Slug)
		}
		for _, cron := range crons.Crons {
			cmd.Stdout.Printf("Cron builds %s on %s \n", cron.Interval, cron.Branch.Name)
		}
	}

	return cli.Success
}

func getCronsForRepo(repo Repository, cRes chan *client.Response, cErr chan error, cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"repository.slug": repo.Slug,
		"include":         "repository.default_branch,branch.last_build",
	}
	res, err := env.Client.PerformAction("crons", "for_repository", params)
	cRes <- res
	cErr <- err
}

type cronsArray []Crons

func (c cronsArray) Len() int {
	return len(c)
}
func (c cronsArray) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c cronsArray) Less(i, j int) bool {
	return c[i].Crons[0].Repository.Slug > c[j].Crons[0].Repository.Slug
}
