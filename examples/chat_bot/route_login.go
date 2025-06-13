package main

import (
	"net/http"

	"github.com/gildas/go-core"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
)

// LoggedInHandler is called after the token is sent back to the app by GCloud CX
func LoggedInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Child("route", "logged_in")
		appConfig, _ := AppConfigFromContext(r.Context())

		client, err := gcloudcx.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the GCloud CX Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client.Organization, _, _ = client.GetMyOrganization(r.Context())

		if err = appConfig.Initialize(r.Context(), client); err != nil {
			log.Errorf("Failed to Initialize App Config", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
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
