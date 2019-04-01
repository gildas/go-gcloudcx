package conversation

import (
	"encoding/json"
	"fmt"
	"time"

	purecloud "github.com/gildas/go-purecloud"
)

// SendMessage sends a message as the chat guest
func (conversation *Conversation) SendMessage(text string) (err error) {
	return conversation.sendBody("standard", text)
}

// SendNotice sends a notice as the chat guest
func (conversation *Conversation) SendNotice(text string) (err error) {
	return conversation.sendBody("notice", text)
}

type sendTypingResponse struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Conversation Conversation `json:"conversation,omitempty"`
	Sender       Member       `json:"sender,omitempty"`
	Timestamp    time.Time    `json:"timestamp,omitempty"`
}

// SendTyping sends a typing indicator to PureCloud as the chat guest
func (conversation *Conversation) SendTyping() (err error) {
	response := &sendTypingResponse{}
	if err = conversation.Client.Post(fmt.Sprintf("webchat/guest/conversations/%s/members/%s/typing", conversation.ID, conversation.Guest.ID), nil, &response, purecloud.RequestOptions{Authorization: "bearer " + conversation.JWT}); err != nil {
		return err
	}
	conversation.Client.Logger.Record("scope", "sendtyping").Infof("Sent successfuly. Response: %+v", response)
	return nil
}

type sendBodyPayload struct {
	BodyType string `json:"bodyType"`
	Body     string `json:"body"`
}

type sendBodyResponse struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Conversation Conversation `json:"conversation,omitempty"`
	Sender       Member       `json:"sender,omitempty"`
	Body         string       `json:"body,omitempty"`
	BodyType     string       `json:"bodyType,omitempty"`
	Timestamp    time.Time    `json:"timestamp,omitempty"`
	SelfURI      string       `json:"selfUri,omitempty"`
}

// sendBody sends a body message as the chat guest
func (conversation *Conversation) sendBody(bodyType, body string) (err error) {
	payload, err := json.Marshal(sendBodyPayload{BodyType: bodyType, Body: body})
	if err != nil {
		return err
	}

	response := &sendBodyResponse{}
	if err = conversation.Client.Post(fmt.Sprintf("webchat/guest/conversations/%s/members/%s/messages", conversation.ID, conversation.Guest.ID), payload, &response, purecloud.RequestOptions{Authorization: "bearer " + conversation.JWT}); err != nil {
		return err
	}
	conversation.Client.Logger.Record("scope", "sendbody").Infof("Sent successfuly. Response: %+v", response)
	return nil
}
