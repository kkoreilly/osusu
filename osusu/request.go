package osusu

import (
	"crypto/sha256"
	"io"
	"net/http"
)

// API access constants (obviously not secure because someone can just copy, but will make it more secure later when actually necessary)
const (
	APIUsername = "z3J8i6gVMyA!H$Ukpvqt5xLos5FgicTeWYf*MtfFU48HMUeMTaMCN59biD^3VxBup@^n7wnWgzCg442!95R9QHnt^6uKZ7f5ip2ycUjbfQ3sWzCZWVP8xgw!dZTn!trD"
	APIPassword = "gbx5T3*UJSALdxAES$n@w2m6b4o949XKMHsApk@Zt4&q3cf$37Jvf#g4#nd95hSnc4K%#h!JD9ifSkDhQyPMT@brtuU!cFxBJwny!ukC$s^ZVPdPzkJm8DvX4bK7to7d"
)

// NewRequest creates an HTTP request to the given URL with basic auth added
func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// just use simple hash on API username and password, not very secure but will make more secure later
	usernameHash := sha256.Sum256([]byte(APIUsername))
	passwordHash := sha256.Sum256([]byte(APIPassword))
	req.SetBasicAuth(string(usernameHash[:]), string(passwordHash[:]))
	return req, nil
}
