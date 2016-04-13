package client

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

func loadToken() string {
	tokenFilePath := tokenFilePath()
	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		createTokenFile()
	}
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		color.Red("Error: Could not read token file!")
		return "error"
	}
	return string(token[:])
}

func createTokenFile() {
	tokenFilePath := tokenFilePath()
	os.Create(tokenFilePath)
	token := promptForToken()
	ioutil.WriteFile(tokenFilePath, []byte(token), 0x644)
}

func promptForToken() string {
	token := ""
	println("I need your travis access token to log you in. Please paste it here:")
	fmt.Scanln(&token)
	return token
}

func tokenFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		color.Red("Error: Could not find home directory!")
		return "~/.travis/token.txt"
	}
	return home + "/.travis/token.txt"
}
