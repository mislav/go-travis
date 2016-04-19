package commands

import (
	"os"
	"strconv"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("show", "displays a build or a job", showCmd) //TODO a job?
}

func showCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": os.Getenv("TRAVIS_REPO"),
		// "include":         "jobs.job",
		// "include":         "jobs.job.Number",
		"limit":   "1",
		"sort_by": "id:desc",
	}

	res, err := client.Travis().PerformAction("builds", "find", params)
	if err != nil {
		color.Red("Build not found.")
		return
	}
	if res.StatusCode > 299 {
		color.Red("Unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) > 0 {
		printCompleteBuild(builds.Builds[0])
	} else {
		color.Red("This repository has no builds yet.")
	}
}

func printCompleteBuild(build Build) {
	y := color.New(color.FgYellow).PrintfFunc()
	c := color.New(color.FgRed, color.Bold).PrintfFunc()
	if build.State == "passed" {
		c = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}

	println("Build #" + build.Number + ":  " + build.Commit.Message)
	y("%-"+strconv.Itoa(12)+"s", "State:")
	c(build.State + "\n")
	printAttribute("Type", build.EventType)
	printAttribute("Branch", build.Branch.Name)
	printAttribute("Duration", formatDuration(build.Duration))
	printAttribute("Started", build.StartedAt)
	printAttribute("Finished", build.StartedAt)
	// println(len(build.Jobs.Jobs))	TODO Jobs still have to be implemented
	// for _, job := range build.Jobs.Jobs {
	// 	printJob(job)
	// }
}

func printAttribute(name string, val string) {
	y := color.New(color.FgYellow).PrintfFunc()
	format := "%-" + strconv.Itoa(12) + "s"
	name += ":"
	y(format, name)
	println(val)
}

func printJob(job Job) {
	println(job.Number)
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
