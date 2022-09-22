package gcloudcx

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// URI represents the Path of a URL (used in SelfURI or in requests, for example)
type URI string

// NewURI creates a new URI with eventually some parameters
//
// path should be a formatter string
func NewURI(path string, args ...interface{}) URI {
	return URI(fmt.Sprintf(path, args...))
}

// WithQuery returns a new URI with the given query
func (uri URI) WithQuery(query Query) URI {
	if len(query) == 0 {
		return uri
	}
	if uri.HasQuery() {
		return URI(fmt.Sprintf("%s&%s", uri, query.Encode()))
	}
	return URI(fmt.Sprintf("%s?%s", uri, query.Encode()))
}

// HasProtocol tells if the URI starts with a protocol/scheme
func (uri URI) HasProtocol() bool {
	matched, _ := regexp.Match(`^[a-z0-9_]+:.*`, []byte(uri))
	return matched
}

// HasPrefix tells if the URI starts with the given prefix
func (uri URI) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(uri), prefix)
}

// HasQuery tells if the URI has a query string
func (uri URI) HasQuery() bool {
	return strings.Contains(string(uri), "?")
}

// Join joins the given paths to the URI
//
// Caveat: does not work if the original URI has a query string
func (uri URI) Join(uris ...URI) URI {
	paths := []string{uri.String()}
	for _, u := range uris {
		paths = append(paths, u.String())
	}
	return URI(path.Join(paths...))
}

// URL returns the URI as a URL
func (uri URI) URL() (*url.URL, error) {
	return url.Parse(string(uri))
}

// String returns the URI as a string
//
// implements the fmt.Stringer interface
func (uri URI) String() string {
	return string(uri)
}
