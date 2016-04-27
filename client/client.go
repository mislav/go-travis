package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/HPI-BP2015H/go-utils/api"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/HPI-BP2015H/go-utils/pathname"
	"github.com/HPI-BP2015H/go-utils/utils"
)

var ignoredHeaders = []string{
	"access-control-allow-credentials",
	"access-control-allow-origin",
	"access-control-expose-headers",
	"connection",
	"date",
	"server",
	"status",
	"strict-transport-security",
	"via",
}

// Travis is a simplyfied constructor for the TravisClient struct.
func Travis(endpoint, token string, debug bool) Client {
	var logger *os.File
	if debug {
		logger = os.Stderr
	}

	tmpdir := pathname.TempDir().Join("travis")
	return NewClient(endpoint, token, logger, tmpdir.String())
}

// TravisClient implement the Client interface and is used for communication with the Travis API.
type TravisClient struct {
	cacheDir string
	manifest *Manifest
	http     *api.Client
	token    string
}

// NewClient is the default constructor for TravisClient.
// For a simplyfied constructor see Travis(endpoint, token string, debug bool).
func NewClient(endpoint, token string, logger *os.File, cacheDir string) *TravisClient {
	rootURL, _ := url.Parse(endpoint)
	http := api.NewClient(rootURL, func(t *api.Transport) {
		if logger != nil {
			debugStream := cli.NewColoredWriter(logger)
			debugStream.PushColor("magenta")

			t.RequestCallback = func(req *http.Request) {
				debugStream.Cprintf("> %s %C(bold)%s://%s%s%C(reset)\n", req.Method, req.URL.Scheme, req.Host, req.URL.RequestURI())
				if req.Body != nil {
					body, _ := ioutil.ReadAll(req.Body)
					debugStream.Cprintf("> %s\n", body)
					req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				}
			}

			t.ResponseCallback = func(res *http.Response) {
				debugStream.Cprintf("< %C(bold)HTTP %d%C(reset)\n", res.StatusCode)

				for name, values := range res.Header {
					if ignoreHeader(name) {
						continue
					}
					value := strings.Join(values, ",")
					fmt.Fprintf(debugStream, "< %s: %s\n", name, value)
				}
				/* Output recived body
				if res.Body != nil {
					body, _ := ioutil.ReadAll(res.Body)
					debugStream.Cprintf("< %s\n", body)
					res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				}
				*/
			}
		}
	})

	return &TravisClient{
		http:     http,
		cacheDir: cacheDir,
		token:    token,
	}
}

func (c *TravisClient) PerformRequest(method, path string, body io.Reader, configure func(*http.Request)) (*Response, error) {
	res, err := c.http.PerformRequest(method, path, body, func(req *http.Request) {
		req.Header.Set("Travis-API-Version", "3")
		if c.token != "" {
			req.Header.Set("Authorization", "token "+c.token)
		}
		if configure != nil {
			configure(req)
		}
	})

	if err != nil {
		return nil, err
	}
	return &Response{Response: res}, nil
}

func (c *TravisClient) PerformAction(resourceName, actionName string, params map[string]string, body map[string]string) (*Response, error) {
	manifest, err := c.Manifest()
	if err != nil {
		return nil, fmt.Errorf("could not get manifest: %q", err.Error())
	}
	resource := manifest.Resource(resourceName)
	if resource == nil {
		return nil, fmt.Errorf("could not find %q resource", resourceName)
	}

	matchingActions := []ResourceAction{}

	for _, action := range resource.AllActions() {
		if actionName == action.Name {
			matchingActions = append(matchingActions, action)
		}
	}

	if len(matchingActions) == 0 {
		return nil, fmt.Errorf("could not find %q action", actionName)
	}

	var path string
	var method string

	for _, action := range matchingActions {
		path, err = utils.ExpandUriTemplate(action.UriTemplate, params)
		if err == nil {
			method = action.RequestMethod
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	var configure func(*http.Request)
	configure = nil
	bodyReader = nil
	if body != nil {
		jsonString, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonString)
		configure = func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	return c.PerformRequest(method, path, bodyReader, configure)
}

func (c *TravisClient) Manifest() (*Manifest, error) {
	if c.manifest != nil {
		return c.manifest, nil
	}

	var res *Response
	var err error

	cache := pathname.NewPathname(c.cacheDir, "manifest.json")
	if cache.Exists() {
		file, err := os.Open(cache.String())
		if err != nil {
			return nil, err
		}
		res = &Response{
			Response: &http.Response{Body: file},
		}
	} else {
		res, err = c.PerformRequest("GET", "/", nil, nil)
		if err != nil {
			return nil, err
		}

		cacheFile, err := cache.Create()
		if err != nil {
			return nil, err
		}

		res.Body = utils.ClosingTeeReader(res.Body, cacheFile)
	}

	c.manifest = &Manifest{}
	err = res.Unmarshal(c.manifest)
	if err != nil {
		return nil, err
	}

	return c.manifest, nil
}

// Token is getter for TravisClient.token.
func (c *TravisClient) Token() string {
	return c.token
}

// SetToken is setter of TravisClient.token
func (c *TravisClient) SetToken(token string) {
	c.token = token
}

type Response struct {
	*http.Response
}

func (r *Response) Unmarshal(dest interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dest)
}

type Manifest struct {
	Config    *ManifestConfig     `json:"config"`
	Resources map[string]Resource `json:"resources"`
}

func (m *Manifest) GithubScopes() []string {
	return m.Config.GithubConfig.Scopes
}

func (m *Manifest) AllResources() []Resource {
	result := []Resource{}
	for name, resource := range m.Resources {
		resource.Name = name
		result = append(result, resource)
	}
	return result
}

func (m *Manifest) Resource(target string) *Resource {
	for name, resource := range m.Resources {
		if name == target {
			resource.Name = name
			return &resource
		}
	}
	return nil
}

type ManifestConfig struct {
	GithubConfig *GithubConfig `json:"github"`
}

type GithubConfig struct {
	Scopes []string `json:"scopes"`
}

type Resource struct {
	Name string
	// Actions map[string][]ResourceAction `json:"actions"`
	Actions     map[string]interface{} `json:"actions"`
	Attributes  []string               `json:"attributes"`
	SortableBy  []string               `json:"sortable_by"`
	DefaultSort string                 `json:"default_sort"`
}

func (r *Resource) AllActions() []ResourceAction {
	result := []ResourceAction{}
	for name, actions := range r.Actions {
		switch a := actions.(type) {
		case []interface{}:
			for _, action := range a {
				action := action.(map[string]interface{})
				method := action["request_method"].(string)
				template := action["uri_template"].(string)
				result = append(result, ResourceAction{
					Name:          name,
					RequestMethod: method,
					UriTemplate:   template,
				})
			}
		}
	}
	return result
}

type ResourceAction struct {
	Name          string
	RequestMethod string `json:"request_method"`
	UriTemplate   string `json:"uri_template"`
}

func ignoreHeader(name string) bool {
	name = strings.ToLower(name)
	for _, ignored := range ignoredHeaders {
		if name == ignored {
			return true
		}
	}
	return false
}
