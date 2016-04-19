package commands

import "github.com/HPI-BP2015H/go-utils/cli"

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "whatsup",
			Help:     "lists most recent builds",
			Function: whatsupCmd,
		},
	)
}

func whatsupCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"include": "branch.last_build",
	}
	repositories, err := GetAllRepositories(params)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get Repositories.")
		cmd.Exit(1)
	}
	for _, repo := range repositories.Repositories {
		if repo.Active && (repo.DefaultBranch.LastBuild != nil) {
			printRepoStatus(repo, cmd)
		}
	}
	cmd.Exit(0)
}

func printRepoStatus(repo Repository, cmd *cli.Cmd) {
	build := repo.DefaultBranch.LastBuild
	PushColorAccordingToBuildStatusBold(*build, cmd)
	cmd.Stdout.Printf("%s ", repo.Slug)
	cmd.Stdout.PopColor()
	PushColorAccordingToBuildStatus(*build, cmd)
	cmd.Stdout.Printf(build.State+": #%s \n", build.Number)
	cmd.Stdout.PopColor()
}
