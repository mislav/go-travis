package commands

import (
	"fmt"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("whoami", "outputs the current user", whoamiCmd)
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

func whoamiCmd(cmd *cli.Cmd) {
	user, err := getCurrentUser()
	if err != nil {
		color.Red("Error: Could not get the current user! \n" + err.Error())
		return
	}
	printUser(user)
}

func getCurrentUser() (User, error) {
	user := User{}
	params := map[string]string{}
	res, err := client.Travis().PerformAction("user", "current", params)
	defer res.Body.Close()
	if err != nil {
		return user, fmt.Errorf("Error: Could not get current user! \n%s", err.Error())
	}
	if res.StatusCode > 299 {
		return user, fmt.Errorf("Error: Unexpected HTTP status: %d\n \n%s", res.StatusCode, err.Error())
	}
	res.Unmarshal(&user)
	return user, nil
}

func printUser(user User) {
	g := color.New(color.FgGreen).PrintfFunc()
	gb := color.New(color.FgGreen, color.Bold).PrintfFunc()
	g("You are ")
	gb(user.Name)
	if (user.Name != user.Login) && (user.Login != "") {
		color.Green(" (%s).", user.Login)
	} else {
		color.Green(".")
	}
}
