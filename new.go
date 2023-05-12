package gcloudcx

import (
	"context"

	"github.com/gildas/go-core"
)

// New create a resource from the Genesys Cloud API
//
// # The object must implement the Initializable interface
//
// Resources can be fetched by their ID:
//
//	user, err := Fetch[gcloudcx.User](context, client, uuid.UUID)
//
//	user, err := Fetch[gcloudcx.User](context, client, gcloudcx.User{ID: uuid.UUID})
//
// or by their URI:
//
//	user, err := Fetch[gcloudcx.User](context, client, gcloudcx.User{}.GetURI(uuid.UUID))
func New[T core.Identifiable, PT interface {
	Initializable
	*T
}](context context.Context, client *Client, parameters ...any) *T {
	id, _, _, log := parseFetchParameters(context, client, parameters...)
	var object T

	PT(&object).Initialize(id, client, log)
	return &object
}
