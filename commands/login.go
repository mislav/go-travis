package commands

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/mislav/go-utils/cli"
	"golang.org/x/oauth2"
)

func init() {
	cli.Register("login", loginCmd)
}

func loginCmd(cmd *cli.Cmd) {
	github := Login()
	user, _, err := github.Users.Get("")
	if err != nil {
		color.Red("Error: Could not connect to Github!")
	}
	color.Green("Success! You are now logged into the account " + *(user.Login) + ".\n")

}

func Login() *github.Client {
	token := client.Travis().Token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}
