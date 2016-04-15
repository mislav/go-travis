package commands

import (
	"io"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("enable", enableCmd)
}

func enableCmd(cmd *cli.Cmd) {
	params := map[string]string{
		"repository.slug": config.RepoSlug(),
	}

	res, err := client.Travis().PerformAction("repository", "enable", params)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(cmd.Stdout, res.Body)
}
