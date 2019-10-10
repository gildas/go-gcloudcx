package conversation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gildas/go-core"
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
	if err = conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s/typing", conversation.ID, conversation.Guest.ID),
		&core.RequestOptions{
			Method:        http.MethodPost, // since payload is empty
			Authorization: "bearer " + conversation.JWT,
		},
		&response,
	); err == nil {
		conversation.Client.Logger.Record("scope", "sendtyping").Infof("Sent successfuly. Response: %+v", response)
	}
	return
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
	response := &sendBodyResponse{}
	if err = conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s/messages", conversation.ID, conversation.Guest.ID),
		&core.RequestOptions{
			Authorization: "bearer " + conversation.JWT,
			Payload:       sendBodyPayload{BodyType: bodyType, Body: body},
		},
		&response,
	); err == nil {
		conversation.Client.Logger.Record("scope", "sendbody").Infof("Sent successfuly. Response: %+v", response)
	}
	return
}