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
	"github.com/google/go-github/github"
	"github.com/howeyc/gopass"
	"golang.org/x/oauth2"
)

func init() {
	cmd := cli.Command{
		Name:     "login",
		Info:     "authenticates against the API and stores the token",
		Function: loginCmd,
	}
	cmd.RegisterFlag(
		cli.Flag{
			Short: "-g",
			Long:  "--github-token",
			Ftype: "GITHUBTOKEN",
			Help:  "identify by GitHub token",
		},
		cli.Flag{
			Short: "-u",
			Long:  "--user",
			Ftype: "LOGIN",
			Help:  "user to log in as",
		},
	)
	cli.AppInstance().RegisterCommand(cmd)
}

func loginCmd(cmd *cli.Cmd) cli.ExitValue {
	env := cmd.Env.(config.TravisCommandConfig)
	message := "%C(green)Successfully logged in as %C(boldgreen)%s%C(reset)%C(green)!%C(reset)\n"
	if env.Token == "" {
		var gitHubAuthorization *github.Authorization

		gitHubToken := cmd.Parameters.String("--github-token", "")
		github, err := loginToGitHub(gitHubToken, cmd.Parameters.String("--user", ""), cmd)
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				cmd.Stderr.Println("Error: The given credentials are not valid.")
				return cli.Failure
			}
			cmd.Stderr.Println("Error: Could not connect to GitHub!\n" + err.Error())
			return cli.Failure
		}
		if gitHubToken == "" {
			gitHubAuthorization, err = getGitHubAuthorization(github)
			if err != nil {
				cmd.Stderr.Println("Error:\n" + err.Error())
				return cli.Failure
			}
			gitHubToken = *gitHubAuthorization.Token
		}
		env.Token, err = getTravisTokenFromGitHubToken(gitHubToken)
		if err != nil {
			cmd.Stderr.Println("Error:\n" + err.Error())
			return cli.Failure
		}
		if gitHubAuthorization != nil {
			github.Authorizations.Delete(*gitHubAuthorization.ID)
		}
	} else {
		if env.Token != env.Config.GetTravisTokenForEndpoint(env.Endpoint) {
			// test travis token if a new one should be set
			_, err := CurrentUser(env.Client)
			if err != nil {
				if strings.Contains(err.Error(), "403") {
					cmd.Stderr.Println("Error: The given token is not valid.")
					return cli.Failure
				}
				cmd.Stderr.Println(err.Error())
				return cli.Failure
			}
		} else {
			message = "%C(green)You are currently already logged in as %C(boldgreen)%s%C(reset)%C(green)! To logout run travis logout.%C(reset)\n"
		}
	}
	env.Config.StoreTravisTokenForEndpoint(env.Token, env.Endpoint)
	env.Client.SetToken(env.Token)
	user, err := CurrentUser(env.Client)
	if err != nil {
		cmd.Stderr.Println("Error:\n" + err.Error())
		return cli.Failure
	}
	cmd.Stdout.Cprintf(message, user)
	return cli.Success
}

// loginToGitHub takes a GitHub token to log into GitHub. If an empty string is
// provided, the user will be prompted for username and password.
func loginToGitHub(token, user string, cmd *cli.Cmd) (*github.Client, error) {
	var github *github.Client
	if token == "" {
		userName, password, err := promptForGitHubCredentials(user, cmd)
		if err != nil {
			return nil, err
		}
		github = loginToGitHubWithUsernameAndPassword(userName, password, cmd)
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

func loginToGitHubWithUsernameAndPassword(username string, password string, cmd *cli.Cmd) *github.Client {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
	client := github.NewClient(tp.Client())
	_, _, err := client.Users.Get("")
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		cmd.Stdout.Print("Two-factor authentication code for " + username + ": ")
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

func showGitHubLoginDisclaimer(cmd *cli.Cmd) {
	cmd.Stdout.Print("We need your ")
	cmd.Stdout.Cprint("bold", "GitHub login")
	cmd.Stdout.Print(" to identify you. \nThis information will ")
	cmd.Stdout.Cprint("bold", "not be sent to Travis CI")
	cmd.Stdout.Print(", only to api.github.com.\nThe password will not be displayed.\n\nTry running with ")
	cmd.Stdout.Cprint("yellow", "--github-token")
	cmd.Stdout.Print(" or ")
	cmd.Stdout.Cprint("yellow", "--auto")
	cmd.Stdout.Println(" if you do not want to enter your password anyway.\n ")
}

func promptForGitHubCredentials(userName string, cmd *cli.Cmd) (string, string, error) {
	showGitHubLoginDisclaimer(cmd)
	if userName == "" {
		cmd.Stdout.Print("Username: ")
		fmt.Scan(&userName)
	}
	cmd.Stdout.Print("Password for " + userName + ": ")
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
