package commands

import (
	"io"
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("enable", "enables a project", enableCmd)
}

func enableCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": os.Getenv("TRAVIS_REPO"),
	}

	res, err := client.Travis().PerformAction("repository", "enable", params)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(cmd.Stdout, res.Body)
}
