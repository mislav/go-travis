package commands

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("login", loginCmd)
}

func loginCmd(cmd *cli.Cmd) {
	token := client.Travis().Token
	println("this is the token: " + token)
	/*
		params := map[string]string{
			"repository.slug":  config.RepoSlug(),
			"build.event_type": "push",
			"limit":            "10",
		}

		res, err := client.Travis().PerformAction("builds", "find", params)
		if err != nil {
			panic(err)
		}
		if res.StatusCode > 299 {
			cmd.Stderr.Printf("unexpected HTTP status: %d\n", res.StatusCode)
			cmd.Exit(1)
		}
	*/
}
