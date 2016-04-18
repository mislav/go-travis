package commands

import (
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
	user := getCurrentUser()
	printUserColorful(user)
}

func getCurrentUser() User {
	params := map[string]string{}
	res, err := client.Travis().PerformAction("user", "current", params)
	defer res.Body.Close()
	if err != nil {
		color.Red("Error: Could not get current user! \n" + err.Error() + "\n Answering with User named Error.")
		//color.Yellow("Fall back to asking Github:")
		//whoamiGithub()
		return User{Name: "Error"}
	}

	if res.StatusCode > 299 {
		color.Red("Error: Unexpected HTTP status: %d\n", res.StatusCode)
		return User{Name: "HTTPError"}
		//cmd.Exit(1)
	}
	user := User{}
	res.Unmarshal(&user)
	return user
}

func whoamiGithub() {
	github, _ := LoginToGitHub("", "")
	user, _, err := github.Users.Get("")
	if err != nil {
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	color.Green("You are logged into the account " + *(user.Login) + ".\n")
}

func printUserColorful(user User) {
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
