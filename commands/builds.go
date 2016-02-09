package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mislav/go-utils/api"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("builds", buildsCmd)
}

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

func ignoreHeader(name string) bool {
	name = strings.ToLower(name)
	for _, ignored := range ignoredHeaders {
		if name == ignored {
			return true
		}
	}
	return false
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number string  `json:"number"`
	State  string  `json:"state"`
	Branch *Branch `json:"branch"`
}

type Branch struct {
	Name string `json:"name"`
}

func buildsCmd(args *cli.Args) {
	url, _ := url.Parse("https://api.travis-ci.org")
	travis := api.NewClient(url, func(t *api.Transport) {
		t.RequestCallback = func(req *http.Request) {
			fmt.Fprintf(os.Stderr, "> %s %s://%s%s\n", req.Method, req.URL.Scheme, req.Host, req.URL.RequestURI())
		}

		t.ResponseCallback = func(res *http.Response) {
			fmt.Fprintf(os.Stderr, "< HTTP %d\n", res.StatusCode)
			for name, values := range res.Header {
				if ignoreHeader(name) {
					continue
				}
				value := strings.Join(values, ",")
				fmt.Fprintf(os.Stderr, "< %s: %s\n", name, value)
			}
		}
	})

	path := args.At(0)
	res, err := travis.PerformRequest("GET", path, nil, func(req *http.Request) {
		req.Header.Set("Travis-API-Version", "3")
	})
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	builds := Builds{}
	json.Unmarshal(body, &builds)

	for _, build := range builds.Builds {
		fmt.Printf("#%s: %s (%s)\n", build.Number, build.State, build.Branch.Name)
	}
}
