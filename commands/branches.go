package commands

import (
	"os"
	"sort"
	"strconv"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("branches", "displays the most recent build for each branch", branchesCmd)
}

type Branches struct {
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Name          string      `json:"name"`
	LastBuild     *Build      `json:"last_build"`
	Repository    *Repository `json:"repo"`
	DefaultBranch bool        `json:"default_branch"`
}

type byBuildNumber []Branch

func (b byBuildNumber) Len() int {
	return len(b)
}
func (b byBuildNumber) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b byBuildNumber) Less(i, j int) bool {
	n, _ := strconv.Atoi(b[i].LastBuild.Number)
	m, _ := strconv.Atoi(b[j].LastBuild.Number)
	if b[j].DefaultBranch {
		return m > n
	} else {
		return n > m
	}
}

func branchesCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": os.Getenv("TRAVIS_REPO"),
		"include":         "repository.default_branch",
	}
	res, err := client.Travis().PerformAction("branches", "find", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		cmd.Exit(1)
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
	}
	branches := Branches{}
	res.Unmarshal(&branches)
	sort.Sort(byBuildNumber(branches.Branches))
	var maxLengthName int
	maxLengthNumber := len(branches.Branches[1].LastBuild.Number)
	for _, branch := range branches.Branches {
		if maxLengthName < len(branch.Name) {
			maxLengthName = len(branch.Name)
		}
	}
	format := "%-" + strconv.Itoa(maxLengthName+3) + "s"
	for _, branch := range branches.Branches {
		printBranch(branch, format, maxLengthNumber, cmd)
	}
	cmd.Exit(0)
}

func printBranch(branch Branch, format string, maxLengthNumber int, cmd *cli.Cmd) {
	cmd.Stdout.Cprintf("yellowbold", format, branch.Name)
	PushColorAccordingToBuildStatus(*branch.LastBuild, cmd)
	cmd.Stdout.Println("#%-"+strconv.Itoa(maxLengthNumber)+"s %s", branch.LastBuild.Number, branch.LastBuild.State)
}
