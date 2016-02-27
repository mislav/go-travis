package commands

import (
	"github.com/mislav/go-travis/client"
	"github.com/mislav/go-travis/config"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("builds", buildsCmd)
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number string  `json:"number"`
	State  string  `json:"state"`
	Branch *Branch `json:"branch"`
}

type Branch struct {
	Name string `json:"name"`
}

func buildsCmd(cmd *cli.Cmd) {
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

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		cmd.Stdout.Printf("#%s: %s (%s)\n", build.Number, build.State, build.Branch.Name)
	}
}
