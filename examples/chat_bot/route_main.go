package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

// findParticipant finds a participant after its user id and purpose
func findParticipant(participants []*purecloud.Participant, user *purecloud.User, purpose string) *purecloud.Participant {
	for _, participant := range participants {
		if participant.Purpose == purpose && participant.User != nil && strings.Compare(user.ID, participant.User.ID) == 0 {
			return participant
		}
	}
	return nil
}

// MainHandler is the main webpage. It displays some login info and a WebChat widget
func MainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Child("route", "main")
		appConfig, err := AppConfigFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the Application Configuration", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Initialize data for the Main Page Template
		viewData := struct {
			Region         string
			DeploymentID   string
			OrganizationID string
			AgentQueue     *purecloud.Queue
			BotQueue       *purecloud.Queue
			BotQueueID     string
			User           *purecloud.User
			ChannelID      string
			WebsocketURL   string
			WebRootPath    string
			LoggedIn       bool
		}{
			WebRootPath: appConfig.WebRootPath,
			LoggedIn:    client.IsAuthorized(),
		}

		// We can use the client only if the agent is logged in...
		if viewData.LoggedIn {
			viewData.Region         = client.Region
			viewData.DeploymentID   = client.DeploymentID
			viewData.OrganizationID = client.Organization.ID
			viewData.AgentQueue     = appConfig.AgentQueue
			viewData.BotQueue       = appConfig.BotQueue
			viewData.User           = appConfig.User

			if appConfig.NotificationChannel != nil {
				viewData.ChannelID    = appConfig.NotificationChannel.ID
				viewData.WebsocketURL = appConfig.NotificationChannel.ConnectURL.String()
			} else {
				// We are logged in the client, but message loops are not started, better logging out!
				log.Warnf("Client is logged in but Notification Channel is not operational, logging out")
				appConfig.Reset()
				client.Logout()
				client.DeleteCookie(w)
				viewData.LoggedIn = false
			}
		}

		log.Infof(`Rendering page "page_main"`)
		pageTemplate, err := template.ParseFiles("page_main.html")
		if err != nil {
			log.Errorf(`Failed to parse page "page_main"`, err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		err = pageTemplate.Execute(w, viewData)
		if err != nil {
			log.Errorf(`Failed to render page "page_main"`, err)
		}
	})
}
