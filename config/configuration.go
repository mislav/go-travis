package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Configuration represents a configuration with access tokens and default endpoint.
// The configuration is stored in a file.
type Configuration struct {
	configurationYML configurationYML
	filePath         string
	fileName         string
	cmd              *cli.Cmd
}

type configurationYML struct {
	DefaultEndpoint string                 `yaml:"default_endpoint"`
	Endpoints       map[string]accessToken `yaml:"endpoints"`
}

type accessToken struct {
	AccessToken string `yaml:"access_token"`
}

const configurationFileName = "config.yml"
const configurationFolderName = ".travis"

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
	if c.configurationYML.Endpoints == nil {
		c.configurationYML.Endpoints = make(map[string]accessToken)
	}
	t := new(accessToken)
	t.AccessToken = token
	c.configurationYML.Endpoints[url] = *t
	c.saveConfigurationYML()
}

// DeleteTravisTokenForEndpoint removes the travis access token for the given endpoint
func (c *Configuration) DeleteTravisTokenForEndpoint(url string) {
	delete(c.configurationYML.Endpoints, url)
	c.saveConfigurationYML()
}

func (c *Configuration) loadConfigurationYML() {
	token, err := ioutil.ReadFile(c.filePath + c.fileName)
	if os.IsNotExist(err) {
		_, err = ioutil.ReadDir(c.filePath)
		if err != nil {
			err = os.Mkdir(c.filePath, 0777)
			if err != nil {
				c.cmd.Stderr.Println("Error: could not create '" + c.filePath + "'!" + err.Error())
				return
			}
		}
		c.cmd.Stdout.Cprintln("yellow", "No configuration file found, creating one!")
		c.saveConfigurationYML()
		return
	} else if err != nil {
		c.cmd.Stderr.Println("Error: Could not read configuration file!")
	}
	err = yaml.Unmarshal(token, &c.configurationYML)
	if err != nil {
		c.cmd.Stderr.Println("Error: Could not parse configuration file!")
	}
}

func (c *Configuration) saveConfigurationYML() {
	out, err := yaml.Marshal(c.configurationYML)
	if err != nil {
		c.cmd.Stderr.Println("Error: Could not marshal configuration!")
	}
	err = ioutil.WriteFile(c.filePath+c.fileName, out, 0644)
	if err != nil {
		c.cmd.Stderr.Println("Error: Could not save configuration!")
	}
}

// DefaultConfiguration creates a new configuration from the default file path
func DefaultConfiguration(cmd *cli.Cmd) *Configuration {
	c := new(Configuration)
	c.cmd = cmd
	c.filePath = defaultConfigurationFilePath(cmd)
	c.fileName = configurationFileName
	c.loadConfigurationYML()
	return c
}

func defaultConfigurationFilePath(cmd *cli.Cmd) string {
	home, err := homedir.Dir()
	if err != nil {
		cmd.Stderr.Println("Error: Could not find home directory!")
		return filepath.FromSlash("~/" + configurationFolderName + "/")
	}
	return filepath.FromSlash(home + "/" + configurationFolderName + "/")
}
