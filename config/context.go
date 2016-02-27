package config

import (
	"bufio"
	"net/url"
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
	}

	file, err := os.Open(os.Getenv("HOME") + "/.travis/config.yml")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		atEndpoints := false
		var endpointUrl *url.URL

		for scanner.Scan() {
			line := scanner.Text()
			if line == "endpoints:" {
				atEndpoints = true
				continue
			} else if atEndpoints {
				if strings.HasPrefix(line, "    ") {
					if host == endpointUrl.Host {
						parts := strings.Split(line, " ")
						return parts[len(parts)-1]
					}
				} else if strings.HasPrefix(line, "  ") {
					endpointUrl, _ = url.Parse(strings.TrimSpace(line))
				} else {
					atEndpoints = false
				}
			}
		}
	}

	return ""
}
