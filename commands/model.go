package commands

import "fmt"

type Branches struct {
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Name          string      `json:"name"`
	LastBuild     *Build      `json:"last_build"`
	Repository    *Repository `json:"repo"`
	DefaultBranch bool        `json:"default_branch"`
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number     string  `json:"number"`
	State      string  `json:"state"`
	StartedAt  string  `json:"started_at"`
	FinishedAt string  `json:"finished_at"`
	Duration   int     `json:"duration"`
	EventType  string  `json:"event_type"`
	Branch     *Branch `json:"branch"`
	Commit     *Commit `json:"commit"`
	Jobs       Jobs    `json:"jobs"`
}

func (b *Build) HasPassed() bool {
	return b.State == "passed"
}

func (b *Build) IsNotYetFinished() bool {
	return ((b.State == "created") || (b.State == "started"))
}

type Commit struct {
	Message string `json:"message"`
}

type Jobs struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Number string `json:"number"`
	State  string `json:"state"`
}

type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"login"`
}

type Repositories struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   string  `json:"description"`
	Active        bool    `json:"active"`
	Private       bool    `json:"private"`
	Owner         *Owner  `json:"owner"`
	DefaultBranch *Branch `json:"default_branch"`
}

func (r *Repository) HasDescription() bool {
	return r.Description != ""
}

// User from Travis
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

func (u User) String() string {
	if u.Name != "" && u.Login != u.Name {
		return fmt.Sprintf("%s (%s)", u.Login, u.Name)
	}
	return u.Login
}
