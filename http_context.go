package purecloud

import (
	"context"
	"net/http"

	"github.com/gildas/go-errors"
)

type key int

// ClientContextKey is the key to store Client in context.Context
const ClientContextKey key = iota + 54329

// ToContext stores this Client in the given context
func (client *Client) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, ClientContextKey, client)
}

// ClientFromContext retrieves a Client from a context
func ClientFromContext(context context.Context) (*Client, error) {
	value := context.Value(ClientContextKey)
	if value == nil {
		return nil, errors.ArgumentMissing.With("Client").WithStack()
	}
	if client, ok := value.(*Client); ok {
		return client, nil
	}
	return nil, errors.ArgumentInvalid.With("Client", value).WithStack()
}

// HttpHandler wraps the client into an http Handler
func (client *Client) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("middleware")

			if client.AuthorizationGrant.AccessToken().LoadFromCookie(r, "pcsession").IsValid() {
				log.Infof("PureCloud Token loaded from cookies")
			} else {
				log.Debugf("PureCloud Token not found in cookies")
			}
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}