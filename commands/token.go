package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("token", tokenCmd)
}

func tokenCmd(cmd *cli.Cmd) {
	token := os.Getenv("TRAVIS_TOKEN")
	endpoint := os.Getenv("TRAVIS_ENDPOINT")

	if len(token) > 0 {
		println("Your access token for " + endpoint + " is " + os.Getenv("TRAVIS_TOKEN"))
	} else {
		color.Red("Not logged in for " + endpoint + ", please run travis login")
	}
}
