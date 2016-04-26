package commands

import (
	"io"
	"strings"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/HPI-BP2015H/go-utils/utils"
)

func init() {
	cmd1 := cli.Command{
		Name:     "api",
		Info:     "makes an (authenticated) API call and prints out the raw result",
		Function: apiCmd,
	}
	cmd2 := cli.Command{
		Name:     "raw",
		Info:     "alias for api",
		Function: apiCmd,
	}
	flag := cli.Flag{
		Short: "-i",
		Long:  "--include-headers",
		Ftype: false,
		Help:  "Include the response headers in the output",
	}

	cmd1.RegisterFlag(flag)
	cmd2.RegisterFlag(flag)
	cli.AppInstance().RegisterCommand(cmd1, cmd2)

}

func unrecognizableUnusedArgs(cmd *cli.Cmd, args *cli.Args) bool {
	if args.Length() > 0 {
		cmd.Stderr.Printf("error: unrecognized argument(s) %q\n", args)
		return true
	}
	return false
}

func apiCmd(cmd *cli.Cmd) cli.ExitValue {
	env := cmd.Env.(config.TravisCommandConfig)
	path := ""
	args := cmd.Args
	if args.Length() > 0 {
		path, args = args.Shift()
	}

	if path == "manifest" {
		showResource := ""
		if args.Length() > 0 {
			showResource, args = args.Shift()
		}
		if unrecognizableUnusedArgs(cmd, args) {
			return cli.Failure
		}
		showManifest(cmd, showResource)
		return cli.Success
	} else if path == "" {
		cmd.Stderr.Println("error: missing PATH argument for request")
		return cli.Failure
	} else {
		if unrecognizableUnusedArgs(cmd, args) {
			return cli.Failure
		}
	}

	res, err := env.Client.PerformRequest("GET", path, nil, nil)
	if err != nil {
		cmd.Stderr.Println(err.Error())
		return cli.Failure
	}
	defer res.Body.Close()

	if cmd.Parameters.Bool("--include-headers") {
		cmd.Stdout.Printf("%s %s\r\n", res.Proto, res.Status)
		for name, values := range res.Header {
			value := strings.Join(values, ",")
			cmd.Stdout.Printf("%s: %s\r\n", name, value)
		}
		cmd.Stdout.Print("\r\n")
	}
	if res.StatusCode < 300 {
		io.Copy(cmd.Stdout, res.Body)
	} else {
		io.Copy(cmd.Stderr, res.Body)
		return cli.Failure
	}
	return cli.Success
}

func showManifest(cmd *cli.Cmd, showResource string) {
	env := cmd.Env.(config.TravisCommandConfig)

	manifest, _ := env.Client.Manifest()

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
			return
		}
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
