package main

import (
	"net/http"
	"text/template"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

const widgetJS = `
  if (!window._genesys) window._genesys = {};
  if (!window._gt)      window._gt = [];

  window._genesys.widgets = {
    main: {
      customStylesheetID: "genesys_widgets_custom",
      theme: "dark",
      lang:  "en",
	  preload: [],
	  debug:   true,
    },
    webchat: {
      userData:           {},
      emojis:             true,
      cometD:             { enabled: false },
      uploadsEnabled:     false,
      enableCustomHeader: true,
      autoInvite: {
        enabled:              false,
        timeToInviteSeconds:  5,
        inviteTimeoutSeconds: 30
      },
      chatButton: {
        enabled:          true,
        openDelay:        1000,
        effectDuration:   300,
        hideDuringInvite: true
      },
      transport: {
        type:            "purecloud-v2-sockets",
        dataURL:         "https://api.{{.Region}}",
        deploymentKey:   "{{.DeploymentID}}",
        orgGuid:         "{{.OrganizationID}}",
        interactionData: {
          routing: {
			targetType:    "QUEUE",
            targetAddress: "{{.QueueName}}"
          }
        }
      }
    }
  };
`
//        dataURL:         "https://api.mypurecloud.com.au",
//        deploymentKey:   "b1a956a4-946a-46ef-a7f1-1a6de6402c18",
//        orgGuid:         "ab560fec-a07e-4e8c-bba6-c68539f42b0e",
//            targetAddress: "AltoCloud_Demo_Queue"

func WidgetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Topic("route").Scope("widget")

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		if client.Organization == nil {
			client.Organization, err = client.GetMyOrganization()
			if err != nil {
				log.Errorf("Failed to retrieve my Organization", err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
		}

		log.Infof("Providing PureCloud Config")
		dictionary := struct {
			Region         string
			DeploymentID   string
			OrganizationID string
			QueueName      string
		}{
			Region:         client.Region,
			DeploymentID:   client.DeploymentID,
			OrganizationID: client.Organization.ID,
			QueueName:      Queue.ID,
		}
		// script := template.Must(template.New("script").Parse(widgetJS))
		scriptTemplate, err := template.New("script").Parse(widgetJS)
		if err != nil {
			log.Errorf("Failed to read the template", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		w.Header().Set("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		err = scriptTemplate.Execute(w, dictionary)
		if err != nil {
			log.Errorf("Failed to execute the template", err)
		}
	})
}