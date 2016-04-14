package client

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

func LoadGithubToken() string {
	return LoadTokenFromPath(GithubTokenFilePath())
}

func DeleteGithubTokenFile() {
	tokenFilePath := GithubTokenFilePath()
	os.Remove(tokenFilePath)
}

func ChangeGithubTokenTo(newToken string) {
	changeToken(GithubTokenFilePath(), newToken)
	color.Green("Your Github token has been updated successfully.")
}

func GithubTokenFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		color.Red("Error: Could not find home directory!")
		return "~/.travis/githubToken.txt"
	}
	return home + "/.travis/githubToken.txt"
}

func promptForGithubToken() string {
	token := ""
	println("I need your github access token to log you in. Please paste it here:")
	fmt.Scanln(&token)
	return token
}

func LoadTravisToken() string {
	return LoadTokenFromPath(TravisTokenFilePath())
}

func DeleteTravisTokenFile() {
	tokenFilePath := TravisTokenFilePath()
	os.Remove(tokenFilePath)
}

func ChangeTravisTokenTo(newToken string) {
	changeToken(TravisTokenFilePath(), newToken)
	color.Green("Your Travis token has been updated successfully.")
}

func TravisTokenFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		color.Red("Error: Could not find home directory!")
		return "~/.travis/travisToken.txt"
	}
	return home + "/.travis/travisToken.txt"
}

func LoadTokenFromPath(tokenFilePath string) string {
	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		token := promptForGithubToken()
		ChangeGithubTokenTo(token)
	}
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		color.Red("Error: Could not read token from path " + tokenFilePath)
		return "error"
	}
	return string(token[:])
}

func changeToken(tokenPath string, newToken string) {
	os.Remove(tokenPath)
	os.Create(tokenPath)
	ioutil.WriteFile(tokenPath, []byte(newToken), 0x644)
}
