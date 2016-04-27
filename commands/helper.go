package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

// CurrentUser returns the user currently logged in into Travis
func CurrentUser(client client.Client) (User, error) {
	user := User{}
	res, err := client.PerformAction("user", "current", map[string]string{}, nil)
	if err != nil {
		return user, fmt.Errorf("Error: Could not get current user! \n%s", err.Error())
	}
	if res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return user, err
		}
		return user, fmt.Errorf("Unexpected HTTP status: %d\n%s\n", res.StatusCode, string(body))
	}
	defer res.Body.Close()
	res.Unmarshal(&user)
	return user, nil
}

// GetAllRepositories returns all the repositories (also those not active in travis)
// of the currently logged in user. also takes params
func GetAllRepositories(params map[string]string, client client.Client) (Repositories, error) {
	if params == nil {
		params = map[string]string{}
	}
	repositories := Repositories{}
	res, err := client.PerformAction("repositories", "for_current_user", params, nil)
	defer res.Body.Close()
	if err != nil {
		return repositories, err
	}
	if res.StatusCode > 299 {
		return repositories, fmt.Errorf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
	}
	res.Unmarshal(&repositories)
	return repositories, nil
}

// NotLoggedIn checks if a valid token is stored.
// Prints error message if the user is not logged in
func NotLoggedIn(cmd *cli.Cmd) bool {
	env := cmd.Env.(config.TravisCommandConfig)

	success := len(env.Token) != 0
	if success {
		user, err := CurrentUser(cmd.Env.(config.TravisCommandConfig).Client)
		success = err == nil && user.Name != ""
	}
	if !success {
		cmd.Stderr.Println("You need to be logged in to do this. For this please run travis login.")
		return true
	}
	return false
}

// NoRepoSpecified checks if a repo is specified and prints an error message if not
func NoRepoSpecified(cmd *cli.Cmd) bool {
	env := cmd.Env.(config.TravisCommandConfig)
	if env.Repo == "" {
		cmd.Stderr.Println("Can't figure out GitHub repo name. Ensure you're in the repo directory, or specify the repo name via -r <owner>/<repo>")
		return true
	}
	return false
}

// PushBoldColorAccordingToBuildStatus pushs a bold green for a passed build,
// yellow for a running or red for all other builds to the ColoredWriter
func PushBoldColorAccordingToBuildStatus(build Build, cmd *cli.Cmd) {
	cmd.Stdout.PushColor("bold" + colorForBuildStatus(build))
}

// PushColorAccordingToBuildStatus pushs a green for a passed build,
// yellow for a running or red for all other builds to the ColoredWriter
func PushColorAccordingToBuildStatus(build Build, cmd *cli.Cmd) {
	cmd.Stdout.PushColor(colorForBuildStatus(build))
}

func colorForBuildStatus(build Build) string {
	if build.HasPassed() {
		return "green"
	}
	if build.IsNotYetFinished() {
		return "yellow"
	}
	return "red"
}
