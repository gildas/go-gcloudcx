package gcloudcx

import (
	"context"

	"github.com/gildas/go-core"
)

// Authorizer describes what a grants should do
type Authorizable interface {
	Authorize(context context.Context, client *Client) (string, error) // Authorize a client with Gcloud
	AccessToken() *AccessToken                                         // Get the Access Token obtained by the Authorizer
	core.Identifiable                                                  // Implements core.Identifiable
}
