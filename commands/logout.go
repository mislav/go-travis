package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
)

func init() {
	cli.Register("logout", "deletes the stored API token", logoutCmd)
}

func logoutCmd(cmd *cli.Cmd) {
	user, _ := getCurrentUser()
	config := config.DefaultConfiguration()
	config.DeleteTravisTokenForEndpoint(os.Getenv("TRAVIS_ENDPOINT"))
	gb := color.New(color.FgGreen, color.Bold).PrintfFunc()
	gb(user.Name)
	color.Green(" is now logged out.")
}
