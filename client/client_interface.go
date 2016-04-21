package client

import (
	"io"
	"net/http"
)

type Client interface {
	PerformRequest(string, string, io.Reader, func(*http.Request)) (*Response, error)
	PerformAction(string, string, map[string]string) (*Response, error)
	Manifest() (*Manifest, error)

	// getters
	Token() string

	// setters
	SetToken(string)
}
