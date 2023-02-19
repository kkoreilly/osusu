package main

import (
	"io"
	"net/http"
)

// NewRequest creates an HTTP request to the given URL with basic auth added
func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(APIUsername, APIPassword)
	return req, nil
}
