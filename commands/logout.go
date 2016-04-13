package commands

import (
	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("logout", logoutCmd)
}

func logoutCmd(cmd *cli.Cmd) {
	client.DeleteTokenFile()
}
