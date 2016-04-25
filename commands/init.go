package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/HPI-BP2015H/go-travis/assets"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cmd := cli.Command{
		Name:     "init",
		Info:     "generates a .travis.yml and enables the project",
		Help:     "%C(bold)Available languages%C(reset): c, clojure, cpp, erlang, go, groovy, haskell, java, node_js, objective-c, perl, php, python, ruby, scala",
		Function: initCmd,
	}
	cmd.RegisterFlag(
		cli.Flag{
			Short: "-f",
			Long:  "--force",
			Ftype: false,
			Help:  "override .travis.yml if it already exists",
		},
		cli.Flag{
			Short: "-k",
			Long:  "--skip-enable",
			Ftype: "LOGIN",
			Help:  "do not enable project, only add .travis.yml",
		},
		cli.Flag{
			Short: "-p",
			Long:  "--print-conf",
			Ftype: false,
			Help:  "print generated config instead of writing to file",
		},
	)
	cli.AppInstance().RegisterCommand(cmd)
}

func initCmd(cmd *cli.Cmd) cli.ExitValue {
	if NoRepoSpecified(cmd) {
		return cli.Failure
	}
	if _, err := os.Stat(".travis.yml"); err == nil && !(cmd.Parameters.IsProvided("--force")) && !(cmd.Parameters.IsProvided("--print-conf")) {
		cmd.Stderr.Println(".travis.yml already exists, use --force to override.")
		return cli.Failure
	}
	var languageName string
	cmd.Stdout.Print("Main programming language used: ")
	fmt.Scanln(&languageName)
	template, err := getYMLTemplateForLanguage(languageName, cmd)
	if err != nil {
		cmd.Stderr.Printf("Could not find a corresponding .travis.yml template for the language %s.\nRun travis help init to see a list of supported languages.\n", languageName)
		return cli.Failure
	}
	if cmd.Parameters.IsProvided("--print-conf") {
		cmd.Stdout.Println(string(template[:]))
		return cli.Success
	}
	err = ioutil.WriteFile(".travis.yml", template, 0644)
	if err != nil {
		cmd.Stderr.Println("Error: Could not save .travis.yml!")
		return cli.Failure
	}
	cmd.Stdout.Cprintln("green", ".travis.yml was created!")
	if !cmd.Parameters.IsProvided("--skip-enable") {
		cmd.Stdout.Println("Enabeling this repository for travis...")
		return enableCmd(cmd)
	}
	return cli.Success
}

func getYMLTemplateForLanguage(languageName string, cmd *cli.Cmd) ([]byte, error) {
	return assets.Asset("init/" + languageName + ".yml")
}
