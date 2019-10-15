package purecloud

import (
	"net/http"
	"github.com/pkg/errors"
	"context"
)

// ClientContextKey is the key to store Client in context.Context
const ClientContextKey = iota

// ToContext stores this Client in the given context
func (client *Client) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, ClientContextKey, client)
}

// ClientFromContext retrieves a Client from a context
func ClientFromContext(context context.Context) (*Client, error) {
	value := context.Value(ClientContextKey)
	if value == nil {
		return nil, errors.New("Context does not contain any purecloud.Client")
	}
	if client, ok := value.(*Client); ok {
		return client, nil
	}
	return nil, errors.New("Invalid purecloud.Client stored in Context")
}

// HttpHandler wraps the client into an http Handler
func (client *Client) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}