package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mislav/go-travis/client"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("api", apiCmd)
}

func apiCmd(args *cli.Args) {
	path := args.Word(0)

	res, err := client.Travis.PerformRequest("GET", path, nil, nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if args.HasFlag("-i") {
		fmt.Printf("%s %s\r\n", res.Proto, res.Status)
		for name, values := range res.Header {
			value := strings.Join(values, ",")
			fmt.Printf("%s: %s\r\n", name, value)
		}
		fmt.Print("\r\n")
	}

	io.Copy(os.Stdout, res.Body)
}
