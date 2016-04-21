# go-travis

Travis API v3 command-line client written in Go.

Goals:

* **Speed.** Go runtime boots extremely fast.
* **Portability.** Precompiled binaries can be distributed without dependencies.
* **Extensibility.** Custom `travis-<foo>` subcommands are usable if in PATH.

Current design:

* Each `commands/*.go` file registers a subcommand.
* A command generally uses `client.Travis()` to perform API actions.
* This Travis HTTP client fetches the API manifest once and performs subsequent
  actions by expanding the URI templates found within.
* Calls to unregistered subcommands are dispatched to `travis-<foo>` executables
  in PATH. The following environment is provided: `TRAVIS_REPO`, `TRAVIS_TOKEN`.
* The custom `travis-<foo>` scripts can be implemented in any scripting language
  and consume the `travis api` subcommand to dispatch manual API requests.

Current supported inputs:

* Global flags: `-r/--repo SLUG`, `-t/--token TOKEN`, `--debug`.
* If repo slug isn't explicitly provided, `git remote` configuration is consulted.
* If token isn't explicitly provided, `~/.travis/config.yml` is consulted.

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

Old travis.rb functionality (checked means it has been implemented in go-travis):

* [ ] `travis accounts`
* [x] `travis branches`
* [ ] `travis cache`
* [ ] `travis cancel`
* [ ] `travis console`
* [x] `travis disable`
* [x] `travis enable`
* [ ] `travis encrypt`
* [ ] `travis encrypt-file`
* [x] `travis endpoint`
* [ ] `travis env`
* [x] `travis help`
* [x] `travis history`
* [ ] `travis init`
* [ ] `travis lint`
* [x] `travis login`
* [x] `travis logout`
* [ ] `travis logs`
* [ ] `travis monitor`
* [ ] `travis open`
* [ ] `travis pubkey`
* [x] `travis raw`
* [ ] `travis report`
* [x] `travis repos`
* [ ] `travis restart`
* [ ] `travis settings`
* [ ] `travis setup`
* [x] `travis show`
* [ ] `travis sshkey`
* [x] `travis status`
* [ ] `travis sync`
* [x] `travis token`
* [x] `travis version`
* [x] `travis whatsup`
* [x] `travis whoami`
