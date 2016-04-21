package commands

import "github.com/HPI-BP2015H/go-utils/cli"

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "version",
			Help:     "outputs the client version",
			Function: versionCmd,
		},
	)
}

func versionCmd(cmd *cli.Cmd) cli.ExitValue {
	app := cli.AppInstance()
	cmd.Stdout.Println(app.Version)
	return cli.Success
}
