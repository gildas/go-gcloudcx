package conversation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gildas/go-logger"
)

// MessageHandlers holds the various callbacks used when receiving messages
type MessageHandlers struct {
	OnClosed       func(conversation *Conversation, message *Message, member *Member)
	OnStateChanged func(conversation *Conversation, message *Message, member *Member)
	OnMessage      func(conversation *Conversation, message *Message, member *Member)
	OnTyping       func(conversation *Conversation, message *Message, member *Member)
}
// HandleMessages is the incoming message loop
func (conversation *Conversation) HandleMessages(handlers MessageHandlers) (err error) {
	if conversation.Socket == nil {
		return fmt.Errorf("Conversation Not Connected")
	}

	log := conversation.Logger.Record("scope", "receive").Child().(*logger.Logger)

	for {
		// get a message body and decode it. (ReadJSON is nice, but in case of unknown message, I cannot get the original string)
		var body []byte

		if _, body, err = conversation.Socket.ReadMessage(); err != nil {
			log.Errorf("Failed to read incoming message", err)
			continue // TODO: Should we bail out?!?
		}

		message := &Message{}
		if err = json.Unmarshal(body, &message); err != nil {
			log.Errorf("Malformed JSON message: %s", body, err)
			continue
		}

		log.Infof("Received: %s (version %s)", message.TopicName, message.Version)
		switch strings.ToLower(message.TopicName) {
		case "channel.metadata":
			if message.EventBody.Message == "WebSocket Heartbeat" {
				log.Debugf("<< %s", message.EventBody.Message)
			} else {
				log.Warnf("Unknown: %s, \n%s,\n%+v", message.TopicName, body, message)
			}

		case "v2.conversations.chats." + conversation.ID + ".members":
			switch strings.ToLower(message.Metadata.Type) {
			case "member-change":
				log.Record("correlation", message.Metadata.CorrelationID).Debugf("Timestamp %s", message.EventBody.Timestamp)
				member, err := conversation.GetMember(message.EventBody.Member.ID)
				if err != nil {
					log.Errorf("Failed to get member info for %s", message.EventBody.Member.ID, err)
					member = &Member{
						ID:    message.EventBody.Member.ID,
						State: message.EventBody.Member.State,
					}
				}
				log.Debugf("%s Member %s (%s) State: %s", member.Role, member.ID, member.DisplayName, member.State)
				if message.EventBody.Member.ID == conversation.Member.ID && message.EventBody.Member.State == "DISCONNECTED" {
					defer conversation.Close()
					if handlers.OnClosed != nil {
						handlers.OnClosed(conversation, message, member)
					}
					return nil // Break the incoming message loop
				}
				if handlers.OnStateChanged != nil {
					handlers.OnStateChanged(conversation, message, member)
				}
			default:
				return fmt.Errorf("Unknown Metadata %s", message.Metadata.Type)
			}

		case "v2.conversations.chats." + conversation.ID + ".messages":
			sender, err := conversation.GetMember(message.EventBody.Sender.ID)
			if err != nil {
				log.Errorf("Failed to get sender info for %s", message.EventBody.Sender.ID, err)
				sender = &Member{ ID: message.EventBody.Sender.ID }
			}
			switch strings.ToLower(message.Metadata.Type) {
			case "message":
				// TODO: Do NOT send the same message twice!
				log.Debugf("Message %s from %s (%s) at %s", message.EventBody.ID, sender.ID, sender.DisplayName, message.EventBody.Timestamp)
				if sender.ID != conversation.Member.ID && handlers.OnMessage != nil {
					handlers.OnMessage(conversation, message, sender)
				}
			case "typing-indicator":
				// TODO: Do NOT send the same message twice!
				log.Debugf("Typing Indicator %s from %s (%s) at %s", message.EventBody.ID, sender.ID, sender.DisplayName, message.EventBody.Timestamp)
				if handlers.OnMessage != nil {
					handlers.OnTyping(conversation, message, sender)
				}
			default:
				log.Warnf("Unknown: %s, \n%s, \n%+v", message.Metadata.Type, body, message)
			}

		default:
			log.Warnf("Unknown: %s, \n%s, \n%+v", message.TopicName, body, message)
		}
	}
}