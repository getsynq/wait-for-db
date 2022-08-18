package util

import (
	"errors"
	nurl "net/url"
	"strings"
)

var errNoScheme = errors.New("no scheme")
var errEmptyURL = errors.New("URL cannot be empty")

// schemeFromURL returns the scheme from a URL string
func SchemeFromURL(url string) (string, error) {
	if url == "" {
		return "", errEmptyURL
	}

	i := strings.Index(url, ":")

	// No : or : is the first character.
	if i < 1 {
		return "", errNoScheme
	}

	return url[0:i], nil
}

// FilterCustomQuery filters all query values starting with `x-`
func FilterCustomQuery(u *nurl.URL) *nurl.URL {
	ux := *u
	vx := make(nurl.Values)
	for k, v := range ux.Query() {
		if len(k) <= 1 || k[0:2] != "x-" {
			vx[k] = v
		}
	}
	ux.RawQuery = vx.Encode()
	return &ux
}
