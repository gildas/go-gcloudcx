package purecloud

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// URI represents the Path of a URL (used in SelfURI, for example)
type URI string

// Addressable describes things that carry a URI (typically /api/v2/things/{{uuid}})
type Addressable interface {
	GetURI() URI
}

// NewURI creates a new URI with eventually some parameters
//
// path should be a formatter string
func NewURI(path string, args ...interface{}) URI {
	return URI(fmt.Sprintf(path, args...))
}

func (uri URI) HasProtocol() bool {
	matched, _ := regexp.Match(`^[a-z0-9_]+:.*`, []byte(uri))
	return matched
}

func (uri URI) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(uri), prefix)
}

func (uri URI) Join(uris ...URI) URI {
	paths := []string{uri.String()}
	for _, u := range uris {
		paths = append(paths, u.String())
	}
	return URI(path.Join(paths...))
}

func (uri URI) URL() (*url.URL, error) {
	return url.Parse(string(uri))
}

func (uri URI) String() string {
	return string(uri)
}
