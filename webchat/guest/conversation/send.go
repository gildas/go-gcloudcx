package conversation

import (
	"fmt"
	"github.com/gildas/go-purecloud"
	"encoding/json"
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
	ID           string       `json:"id,omitifempty"`
	Name         string       `json:"name,omitifempty"`
	Conversation Conversation `json:"conversation,omitifempty"`
	Sender       Member       `json:"sender,omitifempty"`
	Timestamp    string       `json:"timestamp,omitifempty"` // time.Time!?
}

// SendTyping sends a typing indicator to PureCloud as the chat guest
func (conversation *Conversation) SendTyping() (err error) {
	response := &sendTypingResponse{}
	if err = conversation.Client.Post(fmt.Sprintf("webchat/guest/conversations/%s/members/%s/typing", conversation.ID, conversation.Guest.ID), nil, &response, purecloud.RequestOptions{ Authorization: "bearer " + conversation.JWT }); err != nil {
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
	ID           string       `json:"id,omitifempty"`
	Name         string       `json:"name,omitifempty"`
	Conversation Conversation `json:"conversation,omitifempty"`
	Sender       Member       `json:"sender,omitifempty"`
	Body         string       `json:"body,omitifempty"`
	BodyType     string       `json:"bodyType,omitifempty"`
	Timestamp    string       `json:"timestamp,omitifempty"` // time.Time!?
	SelfURI      string       `json:"selfUri,omitifempty"`
}

// sendBody sends a body message as the chat guest
func (conversation *Conversation) sendBody(bodyType, body string) (err error) {
	payload, err := json.Marshal(sendBodyPayload{BodyType: bodyType, Body: body})
	if err != nil {
		return err
	}

	response := &sendBodyResponse{}
	if err = conversation.Client.Post(fmt.Sprintf("webchat/guest/conversations/%s/members/%s/messages", conversation.ID, conversation.Guest.ID), payload, &response, purecloud.RequestOptions{ Authorization: "bearer " + conversation.JWT }); err != nil {
		return err
	}
	conversation.Client.Logger.Record("scope", "sendbody").Infof("Sent successfuly. Response: %+v", response)
	return nil
}
