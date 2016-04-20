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

func showCmd(cmd *cli.Cmd) int {
	env := cmd.Env.(config.TravisCommandConfig)

	params := map[string]string{
		"repository.slug": env.Repo,
		"include":         "job.Number",
		"limit":           "1",
		"sort_by":         "id:desc",
	}

	res, err := env.Client.PerformAction("builds", "find", params)
	if err != nil {
		cmd.Stderr.Println("Build not found.")
		return 1
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d \n", res.StatusCode)
		return 1
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) == 0 {
		cmd.Stderr.Println("This repository has no builds yet.")
		return 1
	}
	printCompleteBuild(builds.Builds[0], cmd)
	return 0
}

func printCompleteBuild(build Build, cmd *cli.Cmd) {

	cmd.Stdout.Cprintf("%C(bold)Build #%s:%C(reset)  %s\n", build.Number, build.Commit.Message)
	cmd.Stdout.Cprintf("%C(yellow)%-"+strconv.Itoa(12)+"s%C(reset)", "State:")
	PushColorAccordingToBuildStatusBold(build, cmd)
	cmd.Stdout.Println(build.State)
	cmd.Stdout.PopColor()
	printAttribute("Type", build.EventType, cmd)
	printAttribute("Branch", build.Branch.Name, cmd)
	printAttribute("Duration", formatDuration(build.Duration), cmd)
	printAttribute("Started", build.StartedAt, cmd)
	printAttribute("Finished", build.StartedAt, cmd)
	//cmd.Stdout.Println(len(build.Jobs.Jobs)) //TODO Jobs still have to be implemented
	// for _, job := range build.Jobs.Jobs {
	// 	printJob(job, cmd)
	// }
}

func printAttribute(name string, val string, cmd *cli.Cmd) {
	format := "%C(yellow)%-" + strconv.Itoa(12) + "s%C(reset)"
	name += ":"
	cmd.Stdout.Cprintf(format, name)
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
