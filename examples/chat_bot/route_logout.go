package main

import (
	"net/http"

	"github.com/gildas/go-logger"
)

// LoggedOutHandler is called after the PureCloud user is logged out
func LoggedOutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Topic("route").Scope("logged_out")
		appConfig, _ := AppConfigFromContext(r.Context())

		if len(appConfig.WebRootPath) > 0 {
			log.Infof("Redirecting to %s", appConfig.WebRootPath)
			http.Redirect(w, r, appConfig.WebRootPath, http.StatusTemporaryRedirect)
		} else {
			log.Infof("Redirecting to /")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	})
}