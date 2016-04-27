package commands

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	cmd := cli.Command{
		Name:     "open",
		Help:     "opens a build or job in the browser",
		Function: openCmd,
	}
	cmd.RegisterFlag(
		cli.Flag{
			Short: "-g",
			Long:  "--github",
			Ftype: false,
			Help:  "Open the corresponding project, compare view or pull request on GitHub",
		},
	)
	cmd.RegisterFlag(
		cli.Flag{
			Short: "-p",
			Long:  "--print",
			Ftype: false,
			Help:  "Print out the URL instead of opening it in a browser",
		},
	)
	cli.AppInstance().RegisterCommand(cmd)
}

func openCmd(cmd *cli.Cmd) cli.ExitValue {
	if NoRepoSpecified(cmd) {
		return cli.Failure
	}
	var url string
	var err error
	number := cmd.Args.Peek(0)
	if numberIsInvalid(number) {
		cmd.Stderr.Println("The given build or job number is invalid.")
		return cli.Failure
	}
	a := strings.Split(number, ".")
	build := a[0]
	job := ""
	if len(a) > 1 {
		job = a[1]
	}
	if cmd.Parameters.IsProvided("--github") {
		url, err = getGithubURLForNumber(build, cmd)
	} else {
		url, err = getTravisURLForNumber(build, job, cmd)
	}
	if err != nil {
		cmd.Stderr.Println("Could not create URL.\n" + err.Error())
		return cli.Failure
	}
	if cmd.Parameters.IsProvided("--print") {

		cmd.Stdout.Cprintf("web view: %C(bold)" + url + "%C(reset)\n")
		return cli.Success
	}
	err = open.Run(url)
	if err != nil {
		cmd.Stderr.Println("Error: " + err.Error())
		return cli.Failure
	}
	return cli.Success
}

func getGithubURLForNumber(build string, cmd *cli.Cmd) (string, error) {
	//TODO
	env := cmd.Env.(config.TravisCommandConfig)
	website := "https://github.com/"
	return website + env.Repo + "#" + build, nil
}

func getTravisURLForNumber(build string, job string, cmd *cli.Cmd) (string, error) {
	env := cmd.Env.(config.TravisCommandConfig)
	website := "https://travis-ci.org/"
	fallbackWebsite := "https://travis-ci.org/404"
	params := map[string]string{
		"repository.slug": env.Repo,
		"build.number":    build,
		"limit":           "1",
	}
	if job != "" {
		params["job.Number"] = job
	}
	res, err := env.Client.PerformAction("builds", "find", params, nil)
	if err != nil {
		cmd.Stderr.Printf("Could not find job or build " + env.Repo + "#" + build + "." + job)
		return fallbackWebsite, errors.New("")
	}
	if res.StatusCode > 299 {
		cmd.Stderr.Printf("Unexpected HTTP status: %d\n", res.StatusCode)
		return fallbackWebsite, errors.New("")
	}

	builds := Builds{}
	res.Unmarshal(&builds)
	if len(builds.Builds) == 0 {
		cmd.Stderr.Printf("Could not find job or build " + env.Repo + "#" + build + "." + job)
		return fallbackWebsite, errors.New("")
	}
	if job != "" && len(builds.Builds[0].Jobs.Jobs) > 0 {
		num, _ := strconv.Atoi(job)
		jobID := strconv.Itoa(builds.Builds[0].Jobs.Jobs[num-1].ID)
		cmd.Stdout.Println(jobID)
		return website + env.Repo + "/jobs/" + jobID, nil
	}
	return website + env.Repo + "/builds/" + strconv.Itoa(builds.Builds[0].ID), nil

}

func numberIsInvalid(number string) bool {
	re := regexp.MustCompile(`^(\d+)(\.(\d+))?$`)
	return !re.MatchString(number)
}
