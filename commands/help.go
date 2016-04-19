package commands

import (
	"sort"
	"strconv"

	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.Register("help", "helps you out when in dire need of information", helpCmd)
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
	cmdNames := cli.CommandNames()
	sort.Sort(byLength(cmdNames))
	maxLength := len(cmdNames[0])
	sort.Strings(cmdNames)
	for _, name := range cmdNames {
		format := "\t%-" + strconv.Itoa(maxLength+3) + "s"
		cmd.Stdout.Printf(format, name)
		cmd.Stdout.Cprintln("yellow", cli.LookupHelp(name))
	}
	cmd.Stdout.Println("\nrun travis help COMMAND for more infos.")
	cmd.Exit(0)
}
