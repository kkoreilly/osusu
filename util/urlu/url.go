// Package urlu (url utilities) provides simple url utility functions
package urlu

import (
	"net/url"
)

// AsURL attempts to convert the given string into a URL.
// If it succeeds, it returns the formatted URL string and true.
// If it fails, it returns an empty string and false.
func AsURL(urlString string) (string, bool) {
	urlObj, err := url.Parse(urlString)
	if err != nil || urlObj.Host == "" {
		return "", false
	}
	urlObj.Scheme = "https"
	return urlObj.String(), true
}
