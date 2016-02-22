package commands

import (
	"io"
	"strings"

	"github.com/mislav/go-travis/client"
	"github.com/mislav/go-utils/cli"
	"github.com/mislav/go-utils/utils"
)

func init() {
	cli.Register("api", apiCmd)
}

func apiCmd(cmd *cli.Cmd) {
	path := cmd.Args.Word(0)

	res, err := client.Travis.PerformRequest("GET", path, nil, nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if cmd.Args.HasFlag("-i") {
		cmd.Stdout.Printf("%s %s\r\n", res.Proto, res.Status)
		for name, values := range res.Header {
			value := strings.Join(values, ",")
			cmd.Stdout.Printf("%s: %s\r\n", name, value)
		}
		cmd.Stdout.Print("\r\n")
	}

	io.Copy(cmd.Stdout, res.Body)
}
