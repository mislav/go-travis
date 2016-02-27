# go-travis

Proof-of-concept Travis API v3 command-line client written in Go.

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

Suggested next steps:

* [ ] Devise a help system for subcommands.
* [ ] Encapsulate error handling such as unrecognized flags or HTTP errors.
* [ ] Enable POSTing data via `travis api` command
* [ ] Add optional line-based output from `travis api` instead of raw JSON
* [ ] `travis login`
* [ ] `travis status`
* [ ] `travis show`
* [ ] `travis restart`
