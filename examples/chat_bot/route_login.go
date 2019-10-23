package main

import (
	"net/http"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

// LoggedInHandler is called after the token is sent back to the app by PureCloud
func LoggedInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Topic("route").Scope("logged_in")
		appConfig, _ := AppConfigFromContext(r.Context())

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client.Organization, _ = client.GetMyOrganization()

		if err = appConfig.Initialize(client); err != nil {
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