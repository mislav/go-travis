package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/howeyc/gopass"
	"golang.org/x/oauth2"
)

func init() {
	cli.Register("login", "authenticates against the API and stores the token", loginCmd)
}

func loginCmd(cmd *cli.Cmd) {
	if os.Getenv("TRAVIS_TOKEN") == "" {
		username, password := promptForGithubCredentials()
		github := loginToGitHubWithUsernameAndPassword(username, password)
		authorization := getGithubAuthorization(github)
		if authorization == nil {
			color.Red("Error: The given credentials/token are not valid. Aborting.")
			return
		}
		token := getTravisToken(*authorization.Token)
		github.Authorizations.Delete(*authorization.ID)
		config := config.DefaultConfiguration()
		config.StoreTravisTokenForEndpoint(token, os.Getenv("TRAVIS_ENDPOINT"))
	}
	whoamiCmd(cmd)
}

func getGithubAuthorization(github *github.Client) *github.Authorization {
	req := createAuthorizationrequest()
	authorization, _, err := github.Authorizations.Create(req)
	if err != nil {
		color.Red(err.Error())
		return authorization
	}
	return authorization
}

func createAuthorizationrequest() *github.AuthorizationRequest {
	note := "Temporary Token"
	req := &github.AuthorizationRequest{
		Note:   &note,
		Scopes: []github.Scope{github.Scope("user"), github.Scope("user:email"), github.Scope("repo")},
	}
	return req
}

func getTravisToken(githubToken string) string {
	type accessToken struct {
		AccessToken string `json:"access_token,omitempty"`
	}
	var token accessToken
	httpClient := http.DefaultClient
	req := travisTokenRequest(githubToken)
	resp, err := httpClient.Do(req)
	bytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		color.Red("Error: Could not get travis token.\n " + err.Error())
		return ""
	}
	return token.AccessToken
}

func travisTokenRequest(githubToken string) *http.Request {
	body := []byte("{\"github_token\":\"" + githubToken + "\"}")
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

//LoginToGithub takes a githubtoken to log into github. If an empty string is
//provided, the user will be prompted for username and password.
func LoginToGithub(token string) *github.Client {
	var github *github.Client
	if token == "" {
		username, password := promptForGithubCredentials()
		github = loginToGitHubWithUsernameAndPassword(username, password)
	} else {
		github = loginToGithubWithToken(token)
	}
	_, _, err := github.Users.Get("")
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			color.Red("Error: The given credentials/token are not valid. \n")
			return nil
		}
		color.Red("Error: Could not connect to github! \n" + err.Error())
		return nil
	}

	return github
}

func promptForGithubCredentials() (string, string) {
	var username string
	showGithubLoginDisclaimer()
	print("Username: ")
	fmt.Scan(&username)
	print("Password for " + username + ": ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		color.Red("Error: could not read password.\n " + err.Error())
		return "", ""
	}
	return username, string(pw)
}

func loginToGithubWithToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func loginToGitHubWithUsernameAndPassword(username string, password string) *github.Client {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
	client := github.NewClient(tp.Client())
	_, _, err := client.Users.Get("")
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		fmt.Print("Two-factor authentication code for " + username + ": ")
		otp, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
	}
	return client
}

func showGithubLoginDisclaimer() {
	y := color.New(color.FgYellow).PrintfFunc()
	b := color.New(color.Bold, color.Underline).PrintfFunc()
	print("We need your ")
	b("GitHub login")
	println(" to identify you.")
	print("This information will ")
	b("not be sent to Travis CI")
	println(", only to api.github.com.")
	println("The password will not be displayed. \n ")
	print("Try running with ")
	y("--github-token")
	print(" or ")
	y("--auto")
	println(" if you don not want to enter your password anyway.\n ")
}
