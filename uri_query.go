package gcloudcx

import (
	"fmt"
	"net/url"
)

// Query represents a query string for URIs
type Query map[string]interface{}

// Encode returns the query as a "URL encoded" string
func (query Query) Encode() string {
	values := url.Values{}
	for key, value := range query {
		values.Set(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
