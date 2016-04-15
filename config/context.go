package config

import (
	"os/exec"
	"regexp"
	"strings"
)

var originNames = []string{"upstream", "github", "origin"}

// RepoSlugFromGit executes git remotes command to get the correct slug of the
// repository in the current folder
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
			remoteURL := matches[0][2]
			matches = githubRe.FindAllStringSubmatch(remoteURL, -1)
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
