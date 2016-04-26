package commands

import (
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {

	cmd1 := cli.Command{
		Name:     "cron",
		Info:     "Shows or modifies cron jobs.",
		Function: listCmd,
	}

	cmd1.RegisterCommand(
		cli.Command{
			Name:     "list",
			Info:     "Lists all cron jobs of a repository.",
			Function: listCmd,
		},
	)

	cmd1.RegisterCommand(
		cli.Command{
			Name:      "delete",
			Info:      "Deletes an existing cron job.",
			Function:  deleteCmd,
			Parameter: "CRON_ID",
		},
	)

	cli.AppInstance().RegisterCommand(cmd1)
}

func listCmd(cmd *cli.Cmd) cli.ExitValue {
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"repository.slug": env.Repo,
	}
	res, err := env.Client.PerformAction("crons", "for_repository", params)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get crons! \n" + err.Error())
		return cli.Failure
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
		return cli.Failure
	}
	crons := Crons{}
	res.Unmarshal(&crons)
	if len(crons.Crons) == 0 {
		cmd.Stdout.Cprintln("yellow", "This repository has no crons!")
		return cli.Success
	}

	for _, cron := range crons.Crons {
		printCron(cron, cmd)
	}

	return cli.Success
}

func deleteCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) {
		return cli.Failure
	}
	cronID := cmd.Args.Peek(0)
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"cron.id": cronID,
	}
	res, err := env.Client.PerformAction("cron", "delete", params)
	if err != nil {
		cmd.Stderr.Println("Error: Request failed! \n" + err.Error())
		return cli.Failure
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
		return cli.Failure
	}
	cron := Cron{}
	res.Unmarshal(&cron)
	cmd.Stdout.Cprintf("Cron with ID %C(boldgreen)%d%C(reset) deleted. \n", cron.ID)
	return cli.Success
}

func printCron(cron Cron, cmd *cli.Cmd) {
	cmd.Stdout.Cprintf("%C(boldgreen)%-18s%C(reset) %d \n", "ID:", cron.ID)
	cmd.Stdout.Cprintf("%C(yellow)%-18s%C(reset) %s \n", "Branch:", cron.Branch.Name)
	cmd.Stdout.Cprintf("%C(yellow)%-18s%C(reset) %s \n", "Interval:", cron.Interval)
	cmd.Stdout.Cprintf("%C(yellow)%-18s%C(reset) %t \n", "Disable by build:", cron.DisableByBuild)
	cmd.Stdout.Cprintf("%C(yellow)%-18s%C(reset) %s \n \n", "Next Enqueuing:", cron.NextEnqueuing)
}
