package purecloud

import (
	"github.com/gorilla/securecookie"
	"github.com/google/uuid"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

var hashKey  = []byte(core.GetEnvAsString("PURECLOUD_SESSION_HASH_KEY", "Pur3Cl0udS3ss10nH@5hK3y"))
var blockKey = []byte(core.GetEnvAsString("PURECLOUD_SESSION_HASH_KEY", "Pur3Cl0udS3ss10nBl0ckK3y"))
var secureCookie = securecookie.New(hashKey, blockKey)

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	ClientID     string                 `json:"clientId"`
	Secret       string                 `json:"clientSecret"`
	RedirectURI  *url.URL               `json:"redirectUri"`
	TokenType    string                 `json:"tokenType"`
	Token        string                 `json:"token"`
	TokenExpires time.Time              `json:"tokenExpires"`
}

var tokenCache = cache.New(cache.NoExpiration, 0)

// Login logs in a Client to PureCloud
//   Uses the credentials stored in the Client
func (client *Client) Login() error {
	return client.LoginWithAuthorizationGrant(client.AuthorizationGrant)
}

// LoginWithAuthorizationGrant logs in a Client to PureCloud with given authorization Grant
func (client *Client) LoginWithAuthorizationGrant(authorizationGrant AuthorizationGrant) (err error) {
	if authorizationGrant == nil {
		return fmt.Errorf("Authorization Grant cannot be nil")
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

			if cookie, err := r.Cookie("pcsession"); err == nil {
				var tokenID string

				log.Debugf("Found Cookie: %s=%s", cookie.Name, cookie.Value)
				if err = secureCookie.Decode("pcsession", cookie.Value, &tokenID); err == nil {
					log.Tracef("Token ID: %s", tokenID)
					if tokenString, ok := tokenCache.Get(tokenID); ok {
						log.Debugf("Found Token String from Cookie: %s", tokenString)
						client.AuthorizationGrant.AccessToken().Token = tokenString.(string)
						next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
						return
					}
					log.Infof("Token Not found, need to login with PureCloud")
				} else {
					log.Warnf("Cookie Value could not be decoded, maybe it was tampered with. Error: %s", err.Error())
				}
			} else {
				log.Infof("Cookie Not Found, need to login with PureCloud, error: %s", err.Error())
			}
			redirectURL, _ := url.Parse("https://login." + client.Region + "/oauth/authorize")
			
			if grant, ok := client.AuthorizationGrant.(*AuthorizationCodeGrant); ok {
				query := redirectURL.Query()
				query.Add("response_type", "code")
				query.Add("client_id",     grant.ClientID)
				query.Add("redirect_uri",  grant.RedirectURL.String())
				redirectURL.RawQuery = query.Encode()
			}
			log.Infof("Redirecting to %s", redirectURL.String())
			http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		})
	}
}

// LoginHandler gets a valid Token from PureCloud using an AuthorizationGrant
func (client *Client) LoginHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("login")
			grant, ok := client.AuthorizationGrant.(*AuthorizationCodeGrant)
			if ! ok {
				log.Errorf("Client's Grant is not an Authorization Code Grant, we cannot continue")
				core.RespondWithError(w, http.StatusInternalServerError, errors.New("Invalid PureCloud OAUTH Grant"))
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

			// Store the Access Token in the Sessions with a UUID
			tokenID := uuid.Must(uuid.NewRandom())
			log.Tracef("Token ID: %s", tokenID)
			tokenCache.SetDefault(tokenID.String(), grant.AccessToken().Token)
			encodedID, _ := secureCookie.Encode("pcsession", tokenID.String())
			log.Tracef("Encoded Token ID: %s", encodedID)
			http.SetCookie(w, &http.Cookie{Name: "pcsession", Value: encodedID, Path: "/", HttpOnly: true})

			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}

// LogoutHandler logs out the current user
func (client *Client) LogoutHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("logout")

			if cookie, err := r.Cookie("pcsession"); err == nil {
				log.Debugf("Found Cookie: %s=%s", cookie.Name, cookie.Value)
				var tokenID string
				if err = secureCookie.Decode("pcsession", cookie.Value, &tokenID); err == nil {
					if tokenString, ok := tokenCache.Get(tokenID); ok {
						log.Debugf("Found Token String from Cookie")
						client.AuthorizationGrant.AccessToken().Token = tokenString.(string)
						if err := client.Logout(); err != nil {
							core.RespondWithError(w, http.StatusInternalServerError, err)
							return
						}
						log.Infof("User is now logged out from PureCloud")
					}
				}
			}
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}