package gcloudcx

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// Authorization contains the login options to connect the client to GCloud
type Authorization struct {
	ClientID     string    `json:"clientId"`
	Secret       string    `json:"clientSecret"`
	RedirectURI  *url.URL  `json:"redirectUri"`
	TokenType    string    `json:"tokenType"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"tokenExpires"`
}

// Login logs in a Client to Gcloud
//
//	Uses the credentials stored in the Client
func (client *Client) Login(context context.Context) (correlationID string, err error) {
	return client.LoginWithAuthorizationGrant(context, client.Grant)
}

// LoginWithAuthorizationGrant logs in a Client to Gcloud with given authorization Grant
func (client *Client) LoginWithAuthorizationGrant(context context.Context, grant Authorizable) (correlationID string, err error) {
	if grant == nil {
		return "", errors.ArgumentMissing.With("Authorization Grant")
	}
	if correlationID, err = grant.Authorize(context, client); err != nil {
		return
	}
	return
}

// AuthorizeHandler validates an incoming Request and sends to Gcloud Authorize process if not
func (client *Client) AuthorizeHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("authorize")

			if client.Grant.AccessToken().LoadFromCookie(r, "pcsession").IsValid() {
				log.Debugf("Found Token from Cookie: %s", client.Grant.AccessToken())
				next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
				return
			}

			log.Infof("Cookie Not Found, need to login with Gcloud CX")
			redirectURL, _ := NewURI("%s/oauth/authorize", client.LoginURL).URL()

			if grant, ok := client.Grant.(*AuthorizationCodeGrant); ok {
				query := redirectURL.Query()
				query.Add("response_type", "code")
				query.Add("client_id", grant.GetID().String())
				query.Add("redirect_uri", grant.RedirectURL.String())
				redirectURL.RawQuery = query.Encode()
			}
			log.Infof("Redirecting to %s", redirectURL.String())
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		})
	}
}

// LoggedInHandler gets a valid Token from GCloud using an AuthorizationGrant
func (client *Client) LoggedInHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("login")
			grant, ok := client.Grant.(*AuthorizationCodeGrant)
			if !ok {
				log.Errorf("Client's Grant is not an Authorization Code Grant, we cannot continue")
				core.RespondWithError(w, http.StatusUnauthorized, errors.ArgumentInvalid.With("grant", "Authorization Code Grant"))
				return
			}

			// Get the Request parameter "code"
			params := r.URL.Query()
			grant.Code = params.Get("code")
			log.Tracef("Authorization Code: %s", grant.Code)
			if correlationID, err := client.Login(r.Context()); err != nil {
				log.Record("gcloudcx-correlation", correlationID).Errorf("Failed to Authorize Grant", err)
				core.RespondWithError(w, http.StatusInternalServerError, err)
				return
			}

			client.Grant.AccessToken().SaveToCookie(w, "pcsession")
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}
