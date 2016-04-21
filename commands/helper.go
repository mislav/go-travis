package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

// CurrentUser returns the user currently logged in into Travis
func CurrentUser(client *client.Client) (User, error) {
	user := User{}
	res, err := client.PerformAction("user", "current", map[string]string{})
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
