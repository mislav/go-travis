# go-travis [![Build Status](https://travis-ci.org/HPI-BP2015H/go-travis.svg?branch=master)](https://travis-ci.org/HPI-BP2015H/go-travis)

Travis API v3 command-line client written in Go.

## Goals:

* **Speed.** Go runtime boots extremely fast.
* **Portability.** Precompiled binaries can be distributed without dependencies.
* **Extensibility.** Custom `travis-<foo>` subcommands are usable if in PATH.

## Current design:

* Each `commands/*.go` file registers a subcommand
* You can register global and command specific flags which than will be parsed and passed automatically
* Each subcommand get a ```TravisCommandConfig```
* This Travis HTTP client fetches the API manifest once and performs subsequent actions by expanding the URI templates found within.
* Calls to unregistered subcommands are dispatched to `travis-<foo>` executables in PATH. The following environment is provided: `TRAVIS_REPO`, `TRAVIS_TOKEN` and `TRAVIS_ENDPOINT`. If the `--debug` flag is provided `TRAVIS_DEBUG` will also be set.
* The custom `travis-<foo>` scripts can be implemented in any scripting language  and may consume the `travis api` subcommand to dispatch manual API requests.

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
* [x] `travis init`
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

### Assets handling:

If you need to change something in the assets folder (e.g. the template .yml files) then afterwards you will need to follow these steps:

* if you have not yet go-bindata installed: `go get -u github.com/jteeuwen/go-bindata/...`
* delete the old assets/assets.go file (!)
* in the go-travis/assets folder run: `go-bindata -o assets.go init` (next to init just append all folders inside the assets folder)
* change the new assets.go file's package from `main` to `assets`
* run `go install`
