package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("status", "checks status of the latest build", statusCmd)
}

func statusCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": os.Getenv("TRAVIS_REPO"),
		"limit":           "1",
		"sort_by":         "id:desc",
	}

	res, err := client.Travis().PerformAction("builds", "find", params)
	if err != nil {
		color.Red("build not found")
		return
	}
	if res.StatusCode > 299 {
		color.Red("unexpected HTTP status: %d\n", res.StatusCode)
		cmd.Exit(1)
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) > 0 {
		printStatus(builds.Builds[0])
	} else {
		color.Red("no builds on this repo")
	}
}

func printStatus(build Build) {

	c := color.New(color.FgRed, color.Bold).PrintfFunc()
	if build.State == "passed" {
		c = color.New(color.FgGreen, color.Bold).PrintfFunc()
	}

	print("build #" + build.Number)
	c(" %s \n", build.State)

}
