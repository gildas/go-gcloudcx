package gcloudcx

import (
	"context"
	"net/http"
)

// Logout logs out a Client from GCloud
func (client *Client) Logout(context context.Context) {
	_ = client.Delete(context, "/tokens/me", nil) // we don't care much about the error as we are logging out
	if client.Grant != nil {
		client.Grant.AccessToken().Reset()
	}
}

// DeleteCookie deletes the GCloud Client cookie from the response writer
func (client *Client) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "pcsession", Value: "", Path: "/", HttpOnly: true, MaxAge: -1, Secure: true})
}

// LogoutHandler logs out the current user
func (client *Client) LogoutHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := client.Logger.Scope("logout")

			if client.Grant.AccessToken().LoadFromCookie(r, "pcsession").IsValid() {
				client.Logout(r.Context())
				client.DeleteCookie(w)
				log.Infof("User is now logged out from GCloud")
			}
			next.ServeHTTP(w, r.WithContext(client.ToContext(r.Context())))
		})
	}
}
