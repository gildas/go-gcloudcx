package purecloud

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	ClientID     string    `json:"clientId"`
	Secret       string    `json:"clientSecret"`
	RedirectURI  *url.URL  `json:"redirectUri"`
	TokenType    string    `json:"tokenType"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"tokenExpires"`
}

// Login logs in a Client to PureCloud
//   Uses the credentials stored in the Client
func (client *Client) Login() error {
	return client.LoginWithAuthorizationGrant(client.AuthorizationGrant)
}

// LoginWithAuthorizationGrant logs in a Client to PureCloud with given authorization Grant
func (client *Client) LoginWithAuthorizationGrant(authorizationGrant AuthorizationGrant) (err error) {
	if authorizationGrant == nil {
		return errors.ArgumentMissing.With("Authorization Grant").WithStack()
	}
	if err = authorizationGrant.Authorize(client); err != nil {
		return err
	}
	return
}

// AuthorizeHandler validates an incoming Request and sends to PureCloud Authorize process if not
func (client *Client) AuthorizeHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("authorize")

			if client.AuthorizationGrant.AccessToken().LoadFromCookie(r, "pcsession").IsValid() {
				log.Debugf("Found Token from Cookie: %s", client.AuthorizationGrant.AccessToken())
				next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
				return
			}

			log.Infof("Cookie Not Found, need to login with PureCloud")
			redirectURL, _ := NewURI("%s/oauth/authorize", client.LoginURL).URL()

			if grant, ok := client.AuthorizationGrant.(*AuthorizationCodeGrant); ok {
				query := redirectURL.Query()
				query.Add("response_type", "code")
				query.Add("client_id", grant.ClientID)
				query.Add("redirect_uri", grant.RedirectURL.String())
				redirectURL.RawQuery = query.Encode()
			}
			log.Infof("Redirecting to %s", redirectURL.String())
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		})
	}
}

// LoggedInHandler gets a valid Token from PureCloud using an AuthorizationGrant
func (client *Client) LoggedInHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("login")
			grant, ok := client.AuthorizationGrant.(*AuthorizationCodeGrant)
			if !ok {
				log.Errorf("Client's Grant is not an Authorization Code Grant, we cannot continue")
				core.RespondWithError(w, http.StatusUnauthorized, errors.New("Invalid PureCloud OAUTH Grant"))
				return
			}

			// Get the Request parameter "code"
			params := r.URL.Query()
			grant.Code = params.Get("code")
			log.Tracef("Authorization Code: %s", grant.Code)
			if err := client.Login(); err != nil {
				log.Errorf("Failed to Authorize Grant", err)
				core.RespondWithError(w, http.StatusInternalServerError, err)
				return
			}

			client.AuthorizationGrant.AccessToken().SaveToCookie(w, "pcsession")
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}
