package purecloud

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gildas/go-core"
)

// ImplicitGrant implements PureCloud's Client Implicit Grants
//   See: https://developer.mypurecloud.com/api/rest/authorization/use-implicit-grant.html
type ImplicitGrant struct {
	ClientID    string
	RedirectURL *url.URL
	Token       AccessToken
}

func (grant *ImplicitGrant) Authorize(client *Client) (err error) {
	log := client.Logger.Scope("authorize").Record("grant", "implicit")

	log.Infof("Authenticating with %s using Client Implicit grant", client.Region)

	// Validates the Grant
	if len(grant.ClientID) == 0   { return fmt.Errorf("Missing Argument ClientID") }
	if grant.RedirectURL   == nil { return fmt.Errorf("Missing Argument RedirectURL") }

	// Resets the token before authenticating
	grant.Token.Reset()
	err = client.SendRequest(
		"https://login." + client.Region + "/oauth/authorize",
		&core.RequestOptions{
			Parameters: map[string]string{
				"response_type": "token",
				"client_id":     grant.ClientID,
				"redirect_uri":  grant.RedirectURL.String(),
			},
		},
		nil, //&response,
	)
	if err != nil { return err }
	return
}

// AccessToken gives the access Token carried by this Grant
func (grant ImplicitGrant) AccessToken() *AccessToken {
	return &grant.Token
}

// ToContext stores this Grant in trhe given context
func (grant *ImplicitGrant) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, GrantContextKey, grant)
}

// HttpHandler waps the grant into an http handler
func (grant *ImplicitGrant) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(grant.ToContext(r.Context())))
		})
	}
}