package conversation

import (
	"fmt"

	"github.com/gildas/go-purecloud"
	"github.com/gorilla/websocket"
)

type createPayload struct {
	OrganizationID string `json:"organizationId"`
	DeploymentID   string `json:"deploymentId"`
	RoutingTarget  Target `json:"routingTarget"`
	Guest          Member `json:"memberInfo"`
}

// Create creates a new chat Conversation in PureCloud
func Create(client *purecloud.Client, organization *purecloud.Organization, target Target, guest Member) (*Conversation, error) {
	var err error

	if organization == nil {
		organization, err = client.GetMyOrganization()
		if err != nil {
			client.Logger.Errorf("Failed to retrieve my Organization", err)
			return nil, err
		}
	}
	payload := createPayload{
		OrganizationID: organization.ID,
		DeploymentID:   client.DeploymentID,
		RoutingTarget:  target,
		Guest:          guest,
	}

	conversation := &Conversation{Client: client, Members: make(map[string]*Member)}
	if err = client.Post("/webchat/guest/conversations", payload, &conversation); err != nil {
		return nil, err
	}
	conversation.Logger = client.Logger.Topic("conversation").Scope("conversation").Record("conversation", conversation.ID).Child()

	conversation.Socket, _, err = websocket.DefaultDialer.Dial(conversation.EventStream, nil)
	if err != nil {
		conversation.Logger.Errorf("Failed while connecting to %s", conversation.EventStream, err)
		conversation.Close()
		return nil, err
	}
	conversation.Guest.AvatarURL   = guest.AvatarURL
	conversation.Guest.DisplayName = guest.DisplayName
	conversation.Guest.Role        = guest.Role
	conversation.Guest.Custom      = guest.Custom
	conversation.Members[conversation.Guest.ID] = &conversation.Guest

	return conversation, nil
}

// Close terminates a conversation (its websocket as well)
func (conversation *Conversation) Close() error {
	log := conversation.Logger.Scope("close").Child()

	if conversation.Socket != nil {
		log.Debugf("Disconnecting websocket")
		if err := conversation.Socket.Close(); err != nil {
			log.Errorf("Failed while closing websocket", err)
			return err
		}
		log.Infof("Websocket disconnected")
	} else if conversation.Client != nil {
		log.Debugf("Guest Member leaving")
		if err := conversation.Client.Delete(fmt.Sprintf("/webchat/guest/conversations/%s/members/%s", conversation.ID, conversation.Guest.ID), nil); err != nil {
			log.Errorf("Failed while guest member was leaving chat", err)
			return err
		}
		log.Infof("Guest Member left")
	}
	conversation.Socket = nil
	conversation.Client = nil
	return nil
}