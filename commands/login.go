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

	"github.com/HPI-BP2015H/go-travis/commands/helper"
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

// loginCmd is currently mising the following args of the original travis cli:
// -T, --auto-token                 try to figure out who you are automatically (might send another apps token to Travis, token will not be stored)
// -p, --auto-password              try to load password from OSX keychain (will not be stored)
// -a, --auto                       shorthand for --auto-token --auto-password
// -M, --no-manual                  do not use interactive login
//     --list-github-token          instead of actually logging in, list found GitHub tokens
//     --skip-token-check           don't verify the token with github
func loginCmd(cmd *cli.Cmd) {
	gitHubTokenFlag, args := cmd.Args.ExtractFlag("-g", "--github-token", "GITHUBTOKEN")
	userFlag, args := args.ExtractFlag("-u", "--user", "LOGIN")
	travisToken := os.Getenv("TRAVIS_TOKEN")
	config := config.DefaultConfiguration()

	if travisToken == "" {
		var gitHubAuthorization *github.Authorization

		gitHubToken := gitHubTokenFlag.String()
		github, err := LoginToGitHub(gitHubToken, userFlag.String())
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				color.Red("Error: The given credentials are not valid.\n")
				return
			}
			color.Red("Error: Could not connect to GitHub!\n" + err.Error())
			return
		}
		if gitHubToken == "" {
			gitHubAuthorization, err = getGitHubAuthorization(github)
			if err != nil {
				color.Red("Error:\n" + err.Error())
				return
			}
			gitHubToken = *gitHubAuthorization.Token
		}
		travisToken, err = getTravisTokenFromGitHubToken(gitHubToken)
		if err != nil {
			color.Red("Error:\n" + err.Error())
			return
		}
		if gitHubAuthorization != nil {
			github.Authorizations.Delete(*gitHubAuthorization.ID)
		}
	} else {
		// TODO: Test Travis Token
	}
	config.StoreTravisTokenForEndpoint(travisToken, os.Getenv("TRAVIS_ENDPOINT"))
	os.Setenv("TRAVIS_TOKEN", travisToken)
	user, err := user.CurrentUser()
	if err != nil {
		color.Red("Error:\n" + err.Error())
		return
	}
	color.Green("Successfully logged in as %s! To logout run travis logout.", user)
}

// LoginToGitHub takes a GitHub token to log into GitHub. If an empty string is
// provided, the user will be prompted for username and password.
func LoginToGitHub(token, user string) (*github.Client, error) {
	var github *github.Client
	if token == "" {
		userName, password, err := promptForGitHubCredentials(user)
		if err != nil {
			return nil, err
		}
		github = loginToGitHubWithUsernameAndPassword(userName, password)
	} else {
		github = loginToGitHubWithToken(token)
	}
	if _, _, err := github.Users.Get(""); err != nil {
		return nil, err
	}
	return github, nil
}

func loginToGitHubWithToken(token string) *github.Client {
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

func getGitHubAuthorization(github *github.Client) (*github.Authorization, error) {
	req := createGitHubAuthorizationRequest()
	authorization, _, err := github.Authorizations.Create(req)
	if err != nil {
		return nil, err
	}
	return authorization, nil
}

func getTravisTokenFromGitHubToken(githubToken string) (string, error) {
	type accessToken struct {
		AccessToken string `json:"access_token,omitempty"`
	}
	var token accessToken
	httpClient := http.DefaultClient
	req, err := createTravisTokenRequest(githubToken)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func createTravisTokenRequest(githubToken string) (*http.Request, error) {
	body := []byte("{\"github_token\":\"" + githubToken + "\"}")
	travisTokenRequest, err := http.NewRequest("POST", "https://api.travis-ci.org/auth/github", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	travisTokenRequest.Header.Add("Accept", "application/vnd.travis-ci.2+json")
	travisTokenRequest.Header.Add("User-Agent", "MyClient/1.0.0")
	travisTokenRequest.Header.Add("Content-Type", "application/json")

	return travisTokenRequest, nil
}

func showGitHubLoginDisclaimer() {
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
	println(" if you do not want to enter your password anyway.\n ")
}

func promptForGitHubCredentials(userName string) (string, string, error) {
	showGitHubLoginDisclaimer()
	if userName == "" {
		print("Username: ")
		fmt.Scan(&userName)
	}
	print("Password for " + userName + ": ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		return "", "", err
	}
	return userName, string(pw), nil
}

func createGitHubAuthorizationRequest() *github.AuthorizationRequest {
	note := "Temporary Token for the Travis CI CLI"
	req := &github.AuthorizationRequest{
		Note:   &note,
		Scopes: []github.Scope{github.Scope("user"), github.Scope("user:email"), github.Scope("repo")},
	}
	return req
}
