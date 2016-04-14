package commands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/mislav/go-utils/cli"
	"golang.org/x/oauth2"
)

func init() {
	cli.Register("login", loginCmd)
}

func loginCmd(cmd *cli.Cmd) {
	github := LoginToGitHub()

	user, _, err := github.Users.Get("")
	if err != nil {
		client.DeleteGithubTokenFile()
		if strings.Contains(err.Error(), "401") {
			color.Red("Error: The given token is not valid. \n")
			return
		}
		color.Red("Error: Could not connect to Github! \n" + err.Error())
		return
	}
	color.Green("Success! You are now logged into the account " + *(user.Login) + ".\n")
	getTravisToken()
}

func LoginToGitHub() *github.Client {
	token := client.Travis().Token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func LoginToTravis() {

}

func getTravisToken() {
	type AccessToken struct {
		AccessToken string `json:"access_token,omitempty"`
	}
	var token AccessToken
	httpClient := http.DefaultClient
	req := travisTokenRequest()
	resp, err := httpClient.Do(req)
	bytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		color.Red("Error in getTravisToken().")
		return
	}
	client.ChangeTravisTokenTo(token.AccessToken)

}
func travisTokenRequest() *http.Request {
	body := []byte("{\"github_token\":\"" + client.LoadGithubToken() + "\"}")
	travisTokenRequest, err := http.NewRequest("POST", "https://api.travis-ci.org/auth/github", bytes.NewBuffer(body))
	if err != nil {
		color.Red("Error: Could not create the travis token request.")
		return nil
	}
	travisTokenRequest.Header.Add("Accept", "application/vnd.travis-ci.2+json")
	travisTokenRequest.Header.Add("User-Agent", "MyClient/1.0.0")
	travisTokenRequest.Header.Add("Content-Type", "application/json")

	return travisTokenRequest
}
