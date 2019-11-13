package main

import (
	"net/http"

	"github.com/gildas/go-logger"
)

// LoggedOutHandler is called after the PureCloud user is logged out
func LoggedOutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Child("route", "logged_out")
		appConfig, _ := AppConfigFromContext(r.Context())

		if err := appConfig.Reset(); err != nil {
			log.Errorf("Failed to reset the app")
		}

		redirectPath := appConfig.WebRootPath
		if len(redirectPath) == 0 {
			redirectPath = "/"
		}
		log.Infof("Redirecting to %s", redirectPath)
		// See: https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#3xx_Redirection
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	})
}