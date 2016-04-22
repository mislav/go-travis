package commands

import (
	"sort"
	"strconv"
	"strings"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "branches",
			Help:     "displays the most recent build for each branch",
			Function: branchesCmd,
		},
	)
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
	}
	return n > m
}

func branchesCmd(cmd *cli.Cmd) cli.ExitValue {
	env := cmd.Env.(config.TravisCommandConfig)
	params := map[string]string{
		"repository.slug": env.Repo,
		"include":         "repository.default_branch,build.commit",
	}
	res, err := env.Client.PerformAction("branches", "find", params)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		return cli.Failure
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		return cli.Failure
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
	format := "%C(yellow)%-" + strconv.Itoa(maxLengthName+3) + "s%C(reset)"
	for _, branch := range branches.Branches {
		if branch.DefaultBranch {
			printBranch(branch, "%C(boldyellow)%-"+strconv.Itoa(maxLengthName+3)+"s%C(reset)", maxLengthNumber, cmd)
		} else {
			printBranch(branch, format, maxLengthNumber, cmd)
		}
	}
	return cli.Success
}

func printBranch(branch Branch, format string, maxLengthNumber int, cmd *cli.Cmd) {
	commitMessage := strings.Replace(branch.LastBuild.Commit.Message, "\n", " ", -1)
	cmd.Stdout.Cprintf(format, branch.Name+":")
	PushColorAccordingToBuildStatus(*branch.LastBuild, cmd)
	cmd.Stdout.Printf("#%-"+strconv.Itoa(maxLengthNumber+1)+"s %s   ", branch.LastBuild.Number, branch.LastBuild.State)
	cmd.Stdout.PopColor()
	cmd.Stdout.Println(commitMessage)
}
