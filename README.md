# go-travis

Travis API v3 command-line client written in Go.

Goals:

* **Speed.** Go runtime boots extremely fast.
* **Portability.** Precompiled binaries can be distributed without dependencies.
* **Extensibility.** Custom `travis-<foo>` subcommands are usable if in PATH.

Current design:

* Each `commands/*.go` file registers a subcommand
* A command get a ```TravisCommandConfig```
* You can register global and command specific flags which than will be parsed and passed automatically
* This Travis HTTP client fetches the API manifest once and performs subsequent
  actions by expanding the URI templates found within.
* Calls to unregistered subcommands are dispatched to `travis-<foo>` executables
  in PATH. The following environment is provided: `TRAVIS_REPO`, `TRAVIS_TOKEN` and `TRAVIS_ENDPOINT`. If the `--debug` flag is provided `TRAVIS_DEBUG` will also be set.
* The custom `travis-<foo>` scripts can be implemented in any scripting language
  and may consume the `travis api` subcommand to dispatch manual API requests.

Compatibility with travis.rb:

* Commands and Flags are the same unless noted below
* Outputs are close to the original one
* The same configuration file is used and the stored tokens and endpoint configuration are taken into account
* A configuration written by go-travis is not in all cases compatible with the old client

Differences in usage compared to travis.rb:

* The `-i, --[no-]interactive` flag has been replaced by the `--no-color` flag

The following global flags exist:
```
  -e, --api-endpoint [URL]       Travis API server to talk to
      --staging                  short-cut for --api-endpoint 'https://api-staging.travis-ci.org/'
  -t, --token [ACCESS_TOKEN]     access token to use
      --debug                    show API requests
  -r, --repo [REPOSITORY_SLUG]   the repository on GitHub
  -h, --help                     show help for the command
      --org                      short-cut for --api-endpoint 'https://api.travis-ci.org/'
      --pro                      short-cut for --api-endpoint 'https://api.travis-ci.com/'
```

The following global flags are still missing:
```
  -i, --[no-]interactive           be interactive and colorful
  -E, --[no-]explode               don't rescue exceptions
      --skip-version-check         don't check if travis client is up to date
      --skip-completion-check      don't check if auto-completion is set up
  -I, --[no-]insecure              do not verify SSL certificate of API endpoint
      --debug-http                 show HTTP(S) exchange
  -X, --enterprise [NAME]          use enterprise setup (optionally takes name for multiple setups)
```

travis.rb functionality that is still missing:

* [ ] `travis accounts`
* [ ] `travis cache`
* [ ] `travis cancel`
* [ ] `travis console`
* [ ] `travis encrypt`
* [ ] `travis encrypt-file`
* [ ] `travis env`
* [ ] `travis init`
* [ ] `travis lint`
* [ ] `travis logs`
* [ ] `travis monitor`
* [ ] `travis open`
* [ ] `travis pubkey`
* [ ] `travis report`
* [ ] `travis restart`
* [ ] `travis settings`
* [ ] `travis setup`
* [ ] `travis sshkey`
* [ ] `travis sync`
