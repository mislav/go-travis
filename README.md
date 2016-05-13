# go-travis [![Build Status](https://travis-ci.org/HPI-BP2015H/go-travis.svg?branch=master)](https://travis-ci.org/HPI-BP2015H/go-travis)

Travis API v3 command-line client written in Go.

## Goals:

* **Speed.** Go runtime boots extremely fast.
* **Portability.** Precompiled binaries can be distributed without dependencies.
* **Extensibility.** Custom `travis-<foo>` subcommands are usable if in PATH.

## Current design:

* Each `commands/*.go` file registers a command
* Commands can have subcommands
* You can register global and (sub)command specific flags which than will be parsed and passed automatically
* A help for all (sub)commands with the registered flags is created automatically
* Each (sub)command gets a ```TravisCommandConfig```
* This Travis HTTP client fetches the API manifest once and performs subsequent actions by expanding the URI templates found within.
* Calls to unregistered commands are dispatched to `travis-<foo>` executables in PATH. The following environment is provided: `TRAVIS_REPO`, `TRAVIS_TOKEN` and `TRAVIS_ENDPOINT`. If the `--debug` flag is provided `TRAVIS_DEBUG` will also be set.
* The custom `travis-<foo>` scripts can be implemented in any scripting language  and may consume the `travis api` command to dispatch manual API requests.

## Compatibility with `travis.rb`:

### General
* Commands and Flags are the same unless noted below
* The output is close to the original
* The same configuration file is used and the stored tokens and endpoint configuration are taken into account
* A configuration written by go-travis is not in all cases compatible with the old client

### Differences in usage compared to `travis.rb`:

* Flags like `-r` are available for all commands (see list below)
* The `raw` command works different
* The `-i, --[no-]interactive` flag has been replaced by the `--no-color` flag
* `--adapter` is no longer available
* For `enable`, `history`, `login`,  `repos`, `show`, `status` and `whatsup` are still some flags missing

### `travis.rb` functionality (checked means it has been implemented in `go-travis`):

* [x] `travis branches` displays the most recent build for each branch
* [ ] `travis cache` lists or deletes repository caches  **Needs v3 Endpoint**
* [ ] `travis cancel` cancels a job or build
* [x] `travis cron` *(not yet merged into the old client)*
* [x] `travis crons` *(not yet merged into the old client)*
* [x] `travis disable` disables a project
* [x] `travis enable` enables a project
* [ ] `travis encrypt` encrypts values for the .travis.yml
* [x] `travis endpoint` displays or changes the API endpoint
* [ ] `travis env` show or modify build environment variables
* [x] `travis help` helps you out when in dire need of information
* [x] `travis history` displays a projects build history
* [x] `travis init` generates a .travis.yml and enables the project
* [ ] `travis lint` display warnings for a .travis.yml
* [x] `travis login` authenticates against the API and stores the token
* [x] `travis logout` deletes the stored API token
* [ ] `travis logs` streams test logs
* [ ] `travis open` opens a build or job in the browser
* [ ] `travis pubkey` prints out a repository's public key  **Needs v3 Endpoint**
* [x] `travis raw` makes an (authenticated) API call and prints out the result
* [ ] `travis report` generates a report useful for filing issues
* [x] `travis repos` lists repositories the user has certain permissions on
* [ ] `travis requests` lists recent requests
* [ ] `travis restart` restarts a build or job
* [ ] `travis settings` access repository settings  **Needs v3 Endpoint**
* [ ] `travis setup` sets up an addon or deploy target
* [x] `travis show` displays a build or job
* [ ] `travis sshkey` checks, updates or deletes an SSH key  **Needs v3 Endpoint**
* [x] `travis status` checks status of the latest build
* [ ] `travis sync` triggers a new sync with GitHub **Needs v3 Endpoint**
* [x] `travis token` outputs the secret API token
* [x] `travis version` outputs the client version
* [x] `travis whatsup` lists most recent builds
* [x] `travis whoami` outputs the current user

### `travis.rb` global flags:

* [x] `-e, --api-endpoint [URL]       Travis API server to talk to`
* [ ] `    --debug-http                 show HTTP(S) exchange`
* [ ] `-E, --[no-]explode               don't rescue exceptions`
* [ ] `-X, --enterprise [NAME]          use enterprise setup (optionally takes name for multiple setups)`
* [ ] `-I, --[no-]insecure              do not verify SSL certificate of API endpoint`
* [x] `    --staging                  short-cut for --api-endpoint 'https://api-staging.travis-ci.org/'`
* [x] `-t, --token [ACCESS_TOKEN]     access token to use`
* [x] `    --debug                    show API requests`
* [x] `-r, --repo [REPOSITORY_SLUG]   the repository on GitHub`
* [x] ` h, --help                     show help for the command`
* [x] `    --org                      short-cut for --api-endpoint 'https://api.travis-ci.org/'`
* [x] `    --pro                      short-cut for --api-endpoint 'https://api.travis-ci.com/'`

### Other missing features:

* automatic updating including global `--skip-version-check` flag
* auto-completion including global `--skip-completion-check` flag
* storing a slug for a folder via `--store-repo`

### `travis.rb` commands which are not going to be implemented in `go-travis`

* `travis accounts` displays accounts and their subscription status
* `travis console` interactive shell
* `travis encrypt-file` encrypts a file and adds decryption steps to .travis.yml
* `travis monitor` live monitor for what's going on

## Assets handling:

If you need to change something in the assets folder (e.g. the template .yml files) then afterwards you will need to follow these steps:

* if you have not yet go-bindata installed: `go get -u github.com/jteeuwen/go-bindata/...`
* delete the old assets/assets.go file (!)
* in the go-travis/assets folder run: `go-bindata -o assets.go init` (next to init just append all folders inside the assets folder)
* change the new assets.go file's package from `main` to `assets`
* run `go install`
