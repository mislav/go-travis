package commands

import (
	"io"
	"strings"

	"github.com/mislav/go-travis/client"
	"github.com/mislav/go-utils/cli"
	"github.com/mislav/go-utils/utils"
)

func init() {
	cli.Register("api", apiCmd)
}

func apiCmd(cmd *cli.Cmd) {
	path := cmd.Args.Word(0)

	if path == "manifest" {
		manifest, _ := client.Travis.Manifest()
		showResource := cmd.Args.Word(1)
		if showResource == "" {
			for _, resource := range manifest.AllResources() {
				actionNames := utils.NewSet()
				for _, action := range resource.AllActions() {
					actionNames.Add(action.Name)
				}
				cmd.Stdout.Print(resource.Name)
				if actionNames.Length() > 0 {
					cmd.Stdout.Printf(": %s", strings.Join(actionNames.Values(), ", "))
				}
				cmd.Stdout.Print("\n")
			}
		} else {
			resource := manifest.Resource(showResource)
			if resource == nil {
				cmd.Stderr.Printf("error: could not find the %q resource\n", showResource)
				cmd.Exit(1)
			} else {
				cmd.Stdout.Cprintf("%C(blue)ATTRIBUTES%C(reset) %v\n", resource.Attributes)
				cmd.Stdout.Cprintf("%C(blue)SORTABLE%C(reset) %v\n", resource.SortableBy)
				if resource.DefaultSort != "" {
					cmd.Stdout.Cprintf("%C(blue)SORTABLE%C(reset) default %1\n", resource.DefaultSort)
				}
				for _, action := range resource.AllActions() {
					cmd.Stdout.Cprintf("%C(blue)ACTION%C(reset) %s %s %s\n", action.Name, action.RequestMethod, action.UriTemplate)
				}
			}
		}
		return
	}

	res, err := client.Travis.PerformRequest("GET", path, nil, nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if cmd.Args.HasFlag("-i") {
		cmd.Stdout.Printf("%s %s\r\n", res.Proto, res.Status)
		for name, values := range res.Header {
			value := strings.Join(values, ",")
			cmd.Stdout.Printf("%s: %s\r\n", name, value)
		}
		cmd.Stdout.Print("\r\n")
	}

	io.Copy(cmd.Stdout, res.Body)
}
