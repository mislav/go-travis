package client

import (
	"io"
	"net/http"
)

// Client is an interface that specifies communication between CLI commands and API.
// It is used to provide different types of clients,
// for example one for communication with the Travis CI servers
// and one for using a fake API for testing purposes.
type Client interface {
	PerformRequest(string, string, io.Reader, func(*http.Request)) (*Response, error)
	PerformAction(string, string, map[string]string) (*Response, error)
	Manifest() (*Manifest, error)

	// getters
	Token() string

	// setters
	SetToken(string)
}
