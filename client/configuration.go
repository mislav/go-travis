package client

import (
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Configuration represents a configuration with access tokens and default endpoint.
// The configuration is stored in a file.
type Configuration struct {
	configurationYML configurationYML
	filePath         string
}

type configurationYML struct {
	DefaultEndpoint string                 `yaml:"default_endpoint"`
	Endpoints       map[string]accessToken `yaml:"endpoints"`
}

type accessToken struct {
	AccessToken string `yaml:"access_token"`
}

const configurationFileName = "config.yml"

// TravisStagingEndpoint is the endpoint of Travis Staging System
const TravisStagingEndpoint = "https://api-staging.travis-ci.org/"

// TravisOrgEndpoint is the endpoint of Travis CI for open source projects
const TravisOrgEndpoint = "https://api.travis-ci.org/"

// TravisProEndpoint is the endpoint of Travis CI for private projects
const TravisProEndpoint = "https://api.travis-ci.com/"

// DeleteDefaultTravisEndpoint deletes the default travis endpoint from the configuration
// and saves the configuration file
func (c *Configuration) DeleteDefaultTravisEndpoint() {
	c.configurationYML.DefaultEndpoint = ""
	c.saveConfigurationYML()
}

// StoreDefaultTravisEndpoint overrides the default travis endpoint in the configuration
// and saves the configuration file
func (c *Configuration) StoreDefaultTravisEndpoint(url string) {
	c.configurationYML.DefaultEndpoint = url
	c.saveConfigurationYML()
}

// GetDefaultTravisEndpoint gets the default travis endpoint from the configuration,
// falls back to TravisOrgEndpoint in case no default is set
func (c *Configuration) GetDefaultTravisEndpoint() string {
	endpoint := c.configurationYML.DefaultEndpoint
	if endpoint != "" {
		return endpoint
	}
	return TravisOrgEndpoint
}

// GetTravisTokenForEndpoint gets the travis access token for the given endpoint
func (c *Configuration) GetTravisTokenForEndpoint(url string) string {
	return c.configurationYML.Endpoints[url].AccessToken
}

// StoreTravisTokenForEndpoint save the given travis access token for the endpoint
// and saves the configuration file
func (c *Configuration) StoreTravisTokenForEndpoint(token, url string) {
	t := new(accessToken)
	t.AccessToken = token
	c.configurationYML.Endpoints[url] = *t
	c.saveConfigurationYML()
}

func (c *Configuration) loadConfigurationYML() {
	token, err := ioutil.ReadFile(c.filePath)
	if os.IsNotExist(err) {
		color.Yellow("Warning: No configuration file found!")
		c.saveConfigurationYML()
		return
	} else if err != nil {
		color.Red("Error: Could not read configuration file!")
	}
	err = yaml.Unmarshal(token, &c.configurationYML)
	if err != nil {
		color.Red("Error: Could not parse configuration file!")
	}
}

func (c *Configuration) saveConfigurationYML() {
	out, err := yaml.Marshal(c.configurationYML)
	if err != nil {
		color.Red("Error: Could not marshall configuration!")
	}
	err = ioutil.WriteFile(c.filePath, out, 0644)
	if err != nil {
		color.Red("Error: Could not save configuration!")
	}
}

// DefaultConfiguration creates a new configuration from the default file path
func DefaultConfiguration() *Configuration {
	c := new(Configuration)
	c.filePath = defaultConfigurationFilePath()
	c.loadConfigurationYML()
	return c
}

func defaultConfigurationFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		color.Red("Error: Could not find home directory!")
		return "~/.travis/" + configurationFileName
	}
	return home + "/.travis/" + configurationFileName
}
