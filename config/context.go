package config

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var originNames = []string{"upstream", "github", "origin"}

func RepoSlugFromGit() string {
	out, err := exec.Command("git", "remote", "-v").CombinedOutput()
	if err != nil {
		return ""
	}

	remoteRe := regexp.MustCompile(`(.+)\s+(.+)\s+\((push|fetch)\)`)
	githubRe := regexp.MustCompile(`[@/](ssh.)?github.com(:/?|/)([^/]+)/([^/]+)`)
	found := map[string]string{}

	for _, line := range strings.Split(string(out), "\n") {
		matches := remoteRe.FindAllStringSubmatch(line, -1)
		if matches != nil {
			remoteName := matches[0][1]
			remoteUrl := matches[0][2]
			matches = githubRe.FindAllStringSubmatch(remoteUrl, -1)
			if matches != nil {
				repoOwner := matches[0][3]
				repoName := strings.TrimSuffix(matches[0][4], ".git")
				found[remoteName] = repoOwner + "/" + repoName
			}
		}
	}

	for _, remoteName := range originNames {
		if slug, set := found[remoteName]; set {
			return slug
		}
	}

	for _, slug := range found {
		return slug
	}

	return ""
}

func RepoSlug() string {
	if envSlug := os.Getenv("TRAVIS_REPO"); envSlug != "" {
		return envSlug
	} else {
		return RepoSlugFromGit()
	}
}

func TokenForHost(host string) string {
	envToken := os.Getenv("TRAVIS_TOKEN")
	if envToken != "" && (host == "api.travis-ci.org" || host == "api.travis-ci.com") {
		return envToken
	} else {
		return ""
	}
}
