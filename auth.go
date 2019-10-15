package purecloud

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// AuthorizationGrant describes the capabilities authorization grants must have
type AuthorizationGrant interface {
	// Authorize this Grant with PureCloud
	Authorize(client *Client) error
	AccessToken() *AccessToken
}

type key int

// GrantContextKey is the key for logger child stored in Context
const GrantContextKey key = iota

// AuthorizationGrantFromContext retrieves a grant from an HTTP Context
func AuthorizationGrantFromContext(context context.Context) (AuthorizationGrant, error) {
	grant := context.Value(GrantContextKey)
	if grant == nil {
		return nil, fmt.Errorf("Context does not contain any grant")
	}
	if implictGrant, ok := grant.(*ImplicitGrant); ok {
		return implictGrant, nil
	}
	return nil, errors.New("Unsupported Authorization Grant")
}