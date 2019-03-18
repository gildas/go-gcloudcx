package conversation

import (
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
	//return conversation.Client.Post("webchat/guest/conversations/"+conversation.ID+"/members/"+conversation.Member.ID+"/messages", payload, &response)
	if err = conversation.Client.Post("webchat/guest/conversations/"+conversation.ID+"/members/"+conversation.Member.ID+"/messages", payload, &response); err != nil {
		return nil
	}
	conversation.Client.Logger.Record("scope", "send").Infof("Sent successfuly. Response: %+v", response)
	return nil
}
