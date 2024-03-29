package main

import (
	"html/template"
	"net/http"

	"github.com/gildas/go-core"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// findParticipant finds a participant after its user id and purpose
func findParticipant(participants []*gcloudcx.Participant, user *gcloudcx.User, purpose string) *gcloudcx.Participant {
	for _, participant := range participants {
		if participant.Purpose == purpose && participant.User != nil && user.ID == participant.User.ID {
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

		client, err := gcloudcx.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the GCloud CX Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Initialize data for the Main Page Template
		viewData := struct {
			Region         string
			DeploymentID   uuid.UUID
			OrganizationID uuid.UUID
			AgentQueue     *gcloudcx.Queue
			BotQueue       *gcloudcx.Queue
			BotQueueID     uuid.UUID
			User           *gcloudcx.User
			ChannelID      uuid.UUID
			WebsocketURL   string
			WebRootPath    string
			LoggedIn       bool
		}{
			WebRootPath: appConfig.WebRootPath,
			LoggedIn:    client.IsAuthorized(),
		}

		// We can use the client only if the agent is logged in...
		if viewData.LoggedIn {
			viewData.Region = client.Region
			viewData.DeploymentID = client.DeploymentID
			viewData.OrganizationID = client.Organization.ID
			viewData.AgentQueue = appConfig.AgentQueue
			viewData.BotQueue = appConfig.BotQueue
			viewData.User = appConfig.User

			if appConfig.NotificationChannel != nil {
				viewData.ChannelID = appConfig.NotificationChannel.ID
				viewData.WebsocketURL = appConfig.NotificationChannel.ConnectURL.String()
			} else {
				// We are logged in the client, but message loops are not started, better logging out!
				log.Warnf("Client is logged in but Notification Channel is not operational, logging out")
				_ = appConfig.Reset(r.Context())
				client.Logout(r.Context())
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
