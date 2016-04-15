package commands

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("logout", "deletes the stored API token", logoutCmd)
}

func logoutCmd(cmd *cli.Cmd) {
	config := client.DefaultConfiguration()
	config.DeleteTravisTokenForEndpoint(os.Getenv("TRAVIS_ENDPOINT"))
}
