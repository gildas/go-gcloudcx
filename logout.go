package purecloud

import (
	"net/http"
)

// Logout logs out a Client from PureCloud
func (client *Client) Logout() {
	_ = client.Delete("/tokens/me", nil) // we don't care much about the error as we are logging out
	if client.AuthorizationGrant != nil {
		client.AuthorizationGrant.AccessToken().Reset()
	}
}

// DeleteCookie deletes the PureCloud Client cookie from the response writer
func (client *Client) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "pcsession", Value: "", Path: "/", HttpOnly: true, MaxAge: -1})
}

// LogoutHandler logs out the current user
func (client *Client) LogoutHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("logout")

			if client.AuthorizationGrant.AccessToken().LoadFromCookie(r, "pcsession").IsValid() {
				client.Logout()
				client.DeleteCookie(w)
				log.Infof("User is now logged out from PureCloud")
			}
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}
