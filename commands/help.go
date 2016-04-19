package commands

import (
	"fmt"
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

type stringByLength []string

func (s stringByLength) Len() int {
	return len(s)
}
func (s stringByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s stringByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

type flagByLong []cli.Flag

func (s flagByLong) Len() int {
	return len(s)
}
func (s flagByLong) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s flagByLong) Less(i, j int) bool {
	return len(s[i].Long) > len(s[j].Long)
}

type flagByLength []cli.Flag

func (s flagByLength) Len() int {
	return len(s)
}
func (s flagByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s flagByLength) Less(i, j int) bool {
	return flagLen(s[i]) > flagLen(s[j])
}

func helpCmd(cmd *cli.Cmd) {
	cmd.Stdout.Println("Usage: travis COMMAND [OPTIONS]\n ")
	cmd.Stdout.Println("Available commands:\n ")
	cmdNames := commandNames()
	sort.Sort(stringByLength(cmdNames))
	maxLength := len(cmdNames[0])
	sort.Strings(cmdNames)
	for _, name := range cmdNames {
		format := "\t%-" + strconv.Itoa(maxLength+3) + "s"
		cmd.Stdout.Printf(format, name)
		cmd.Stdout.Cprintln("yellow", lookUpHelp(name))
	}
	cmd.Stdout.Println("\nAvailable Options:\n ")
	flags := globalOptions()
	sort.Sort(flagByLength(flags))
	maxLength = flagLen(flags[0])
	sort.Sort(flagByLong(flags))
	for _, flag := range flags {
		cmd.Stdout.Print("\t")
		if flag.Short != "" {
			cmd.Stdout.Print(flag.Short + ", ")
		} else {
			cmd.Stdout.Print("    ")
		}
		if flag.Ftype != false {
			output := fmt.Sprintf("%v [%v]", flag.Long, flag.Ftype)
			format := "%-" + strconv.Itoa(maxLength+3) + "s"
			cmd.Stdout.Printf(format, output)
		} else {
			format := "%-" + strconv.Itoa(maxLength+3) + "s"
			cmd.Stdout.Printf(format, flag.Long)
		}
		cmd.Stdout.Cprintln("yellow", flag.Help)
	}
	cmd.Stdout.Println("\nRun travis help COMMAND for more infos.")
	cmd.Exit(0)
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

func globalOptions() []cli.Flag {
	app := cli.AppInstance()
	flags := app.Flags()
	result := make([]cli.Flag, 0, len(flags))
	for _, flag := range flags {
		result = append(result, *flag)
	}
	return result
}

func lookUpHelp(cmdName string) string {
	app := cli.AppInstance()
	cmds := app.Commands()
	return cmds[cmdName].Help
}

func flagLen(flag cli.Flag) int {
	result := len(flag.Long)
	if flag.Ftype != false {
		result += len(fmt.Sprintf(" [%v]", flag.Ftype))
	}
	return result
}
