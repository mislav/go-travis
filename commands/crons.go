package commands

import (
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
	//login
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

	var c chan *client.Response = make(chan *client.Response)
	for _, repo := range repositories.Repositories {
		go getCronsForRepo(repo, c, cmd)
	}

	for range repositories.Repositories {
		res := <-c
		crons := Crons{}
		res.Unmarshal(&crons)

		if len(crons.Crons) != 0 {
			repo := crons.Crons[0].Repository
			if repo.DefaultBranch.LastBuild != nil {
				PushColorAccordingToBuildStatusBold(*repo.DefaultBranch.LastBuild, cmd)
				cmd.Stdout.Println(repo.Slug)
				cmd.Stdout.PopColor()

				for _, cron := range crons.Crons {
					cmd.Stdout.Printf("Cron %d builds %s on %s \n", cron.ID, cron.Interval, cron.Branch.Name)
				}
			}
		}
	}

	return cli.Success
}

func getCronsForRepo(repo Repository, c chan *client.Response, cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"repository.slug": repo.Slug,
		"include":         "repository.default_branch,branch.last_build",
	}
	res, err := env.Client.PerformAction("crons", "for_repository", params)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get crons for " + repo.Slug + err.Error())
		//i might want to exit the crons cmd now ... but i cant
	}
	c <- res
}
