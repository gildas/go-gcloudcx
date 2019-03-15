package conversation

import (
	"encoding/json"

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
	// TODO sanitizing...
	payload, err := json.Marshal(createPayload{
		OrganizationID: client.Organization.ID,
		DeploymentID:   client.DeploymentID,
		RoutingTarget:  target,
		Member:         member,
	})
	if err != nil {
		return nil, err
	}

	conversation := &Conversation{}
	if err = client.Post("webchat/guest/conversations", payload, &conversation); err != nil {
		return nil, err
	}
	return conversation, nil
}