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
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Topic("route").Scope("logged_in")

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client.Organization, _ = client.GetMyOrganization()

		if len(AgentQueue.ID) == 0 {
			// TODO: Code this again and cleanly!
			queueName := AgentQueue.Name
			AgentQueue, err = client.FindQueueByName(queueName)
			if err != nil {
				log.Errorf("Failed to retrieve the PureCloud Queue %s", queueName, err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
			log.Record("queue", AgentQueue).Infof("Agent Queue: %s (%s)", AgentQueue.Name, AgentQueue.ID)
		}
		if len(WebRootPath) > 0 {
			log.Infof("Redirecting to %s", WebRootPath)
			http.Redirect(w, r, WebRootPath, http.StatusTemporaryRedirect)
		} else {
			log.Infof("Redirecting to /")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	})
}

// LoginHandler validates the various variables in the Application
func LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log, _ := logger.FromContext(r.Context())
		log = log.Topic("route").Scope("login")

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client.Organization, _ = client.GetMyOrganization()

		if len(AgentQueue.ID) == 0 {
			// TODO: Code this again and cleanly!
			queueName := AgentQueue.Name
			AgentQueue, err = client.FindQueueByName(queueName)
			if err != nil {
				log.Errorf("Failed to retrieve the PureCloud Queue %s", queueName, err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
			log.Record("queue", AgentQueue).Infof("Agent Queue: %s (%s)", AgentQueue.Name, AgentQueue.ID)
		}
		if len(WebRootPath) > 0 {
			log.Infof("Redirecting to %s", WebRootPath)
			http.Redirect(w, r, WebRootPath, http.StatusTemporaryRedirect)
		} else {
			log.Infof("Redirecting to /")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	})

}