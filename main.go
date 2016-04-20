package main

import (
	"os"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/run"
)

// main the current implementation is not respection the debug flag
// The following arguments from the original travis cli are missing:
// -i, --[no-]interactive           be interactive and colorful
// -E, --[no-]explode               don't rescue exceptions
//     --skip-version-check         don't check if travis client is up to date
//     --skip-completion-check      don't check if auto-completion is set up
// -I, --[no-]insecure              do not verify SSL certificate of API endpoint
//     --debug-http                 show HTTP(S) exchange
// -X, --enterprise [NAME]          use enterprise setup (optionally takes name for multiple setups)
func main() {
	os.Exit(run.Run(client.Travis))
}
