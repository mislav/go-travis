package config

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-utils/cli"
)

// TravisCommandConfig stores additional variables which are passed to each command
type TravisCommandConfig struct {
	Repo     string
	Endpoint string
	Token    string
	Debug    bool
	Config   *Configuration
	Client   *client.Client
	cli.CommandConfig
}
