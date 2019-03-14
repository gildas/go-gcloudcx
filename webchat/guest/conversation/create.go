package conversation

import (
	"encoding/json"

	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

type createPayload struct {
	OrganizationID string     `json:"organizationId"`
	DeploymentID   string     `json:"deploymentId"`
	RoutingTarget  Target     `json:"routingTarget"`
	Member         ChatMember `json:"memberInfo"`
}

// Create creates a new chat Conversation in PureCloud
func Create(client *purecloud.Client, target Target, member ChatMember) (*Conversation, error) {
	log := client.Logger.Record("scope", "create_conversation").Child().(*logger.Logger)

	// sanitizing...

	log.Debugf("Creating a new HTTP Request")
	payload, err := json.Marshal(createPayload{
		OrganizationID: client.Organization.ID,
		DeploymentID:   client.DeploymentID,
		RoutingTarget:  target,
		Member:         member,
	})
	if err != nil {
		log.Errorf("Error while encoding payload", err)
		return nil, err
	}

	conversation := &Conversation{}
	err = client.Post("webchat/guest/conversations", payload, &conversation)

	log.Debugf("Success: %+v", conversation)
	return conversation, nil
}