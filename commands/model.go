package commands

import "fmt"

// Branches is an array with branches from Travis
type Branches struct {
	Branches []Branch `json:"branches"`
}

// Branch from Travis
type Branch struct {
	Name          string      `json:"name"`
	LastBuild     *Build      `json:"last_build"`
	Repository    *Repository `json:"repo"`
	DefaultBranch bool        `json:"default_branch"`
}

// Builds is an array with builds from Travis
type Builds struct {
	Builds []Build `json:"builds"`
}

// Build from Travis
type Build struct {
	ID         int     `json:"id"`
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

// HasPassed checks if a Build has passed
func (b *Build) HasPassed() bool {
	return b.State == "passed"
}

// IsNotYetFinished checks if a Build is in the stae created or started
func (b *Build) IsNotYetFinished() bool {
	return ((b.State == "created") || (b.State == "started"))
}

// Commit from Travis
type Commit struct {
	Message string `json:"message"`
}

// Crons is an array with crons from Travis
type Crons struct {
	Crons []Cron `json:"crons"`
}

// Cron from Travis
type Cron struct {
	ID             int         `json:"id"`
	Repository     *Repository `json:"repository"`
	Branch         *Branch     `json:"branch"`
	Interval       string      `json:"interval"`
	DisableByBuild bool        `json:"disable_by_build"`
	NextEnqueuing  string      `json:"next_enqueuing"`
}

// Jobs is an array with jobs from Travis
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

// Job from Travis
type Job struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	State  string `json:"state"`
}

// Owner from Travis
type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"login"`
}

// Repositories is an array with repositories from Travis
type Repositories struct {
	Repositories []Repository `json:"repositories"`
}

// Repository from Travis
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

// HasDescription checks if the Repository description is not empty
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
