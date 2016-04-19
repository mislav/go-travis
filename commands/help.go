package commands

import (
	"sort"
	"strconv"

	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "help",
			Help:     "helps you out when in dire need of information",
			Function: helpCmd,
		},
	)
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func helpCmd(cmd *cli.Cmd) {
	cmd.Stdout.Println("Usage: travis COMMAND ...\n ")
	cmd.Stdout.Println("Available commands: \n ")
	cmdNames := commandNames()
	sort.Sort(byLength(cmdNames))
	maxLength := len(cmdNames[0])
	sort.Strings(cmdNames)
	for _, name := range cmdNames {
		format := "\t%-" + strconv.Itoa(maxLength+3) + "s"
		cmd.Stdout.Printf(format, name)
		cmd.Stdout.Cprintln("yellow", lookUpHelp(name))
	}
	println("\nrun travis help COMMAND for more infos")
}

func commandNames() []string {
	app := cli.AppInstance()
	cmds := app.Commands()
	result := make([]string, 0, len(cmds))
	for cmdName := range cmds {
		result = append(result, cmdName)
	}
	return result
}

func lookUpHelp(cmdName string) string {
	app := cli.AppInstance()
	cmds := app.Commands()
	return cmds[cmdName].Help
}
