package commands

import (
	"strconv"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "show",
			Help:     "displays a build or a job", //TODO a job?
			Function: showCmd,
		},
	)
}

func showCmd(cmd *cli.Cmd) {
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug": env.Repo,
		// "include":         "jobs.job",
		// "include":         "jobs.job.Number",
		"limit":   "1",
		"sort_by": "id:desc",
	}

	res, err := env.Client.PerformAction("builds", "find", params)
	if err != nil {
		cmd.Stderr.Println("Build not found.")
		cmd.Exit(1)
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d \n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) > 0 {
		printCompleteBuild(builds.Builds[0], cmd)
	} else {
		cmd.Stderr.Println("This repository has no builds yet.")
	}
	cmd.Exit(0)
}

func printCompleteBuild(build Build, cmd *cli.Cmd) {

	cmd.Stdout.Println("Build #" + build.Number + ":  " + build.Commit.Message)
	PushColorAccordingToBuildStatusBold(build, cmd)
	cmd.Stdout.Cprint("yellow", "%-"+strconv.Itoa(12)+"s", "State:")
	cmd.Stdout.Println(build.State)
	cmd.Stdout.PopColor()
	printAttribute("Type", build.EventType, cmd)
	printAttribute("Branch", build.Branch.Name, cmd)
	printAttribute("Duration", formatDuration(build.Duration), cmd)
	printAttribute("Started", build.StartedAt, cmd)
	printAttribute("Finished", build.StartedAt, cmd)
	// cmd.Stdout.Println(len(build.Jobs.Jobs))	TODO Jobs still have to be implemented
	// for _, job := range build.Jobs.Jobs {
	// 	printJob(job, cmd)
	// }
}

func printAttribute(name string, val string, cmd *cli.Cmd) {
	format := "%-" + strconv.Itoa(12) + "s"
	name += ":"
	cmd.Stdout.Cprint("yellow", format, name)
	cmd.Stdout.Println(val)
}

func printJob(job Job, cmd *cli.Cmd) {
	cmd.Stdout.Println(job.Number)
}

func formatDuration(duration int) string {
	var res string
	minutes := int(duration / 60)
	seconds := int(duration - (minutes * 60))
	if minutes > 0 {
		res += strconv.Itoa(minutes) + " min "
	}
	res += strconv.Itoa(seconds) + " sec"
	return res
}
