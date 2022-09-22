package main

import (
	"context"
	"strings"
	"time"

	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-request"
)

// MessageLoop receives PureCloud Notification Topics and handles them
func MessageLoop(config *AppConfig, client *gcloudcx.Client) {
	log := config.Logger.Child("topic", "process")

	channel := config.NotificationChannel
	context := context.Background()

	// Processing Received NotificationTopic by reading the chan
	// We do this in a non-blocking way with a timeout to not loop too fast
	// TODO: Add a chan to stop the goroutine that would be written when the route /logout is executed
	for {
		select {
		case receivedTopic := <-channel.TopicReceived:
			if receivedTopic == nil {
				log.Infof("Terminating Topic message loop")
				return
			}
			log.Debugf("Received topic: %s", receivedTopic)
			switch topic := receivedTopic.(type) {
			case *gcloudcx.UserConversationChatTopic:
				log = log.Records("user", topic.User.ID, "conversation", topic.Conversation.ID)
				log.Infof("User %s, Conversation: %s (state: %s)", topic.User, topic.Conversation, topic.Conversation.State)
				for i, participant := range topic.Participants {
					log.Infof("  Participant #%d: id=%s, name=%s, purpose=%s, state=%s", i, participant.ID, participant.Name, participant.Purpose, participant.State)
				}
				participant := findParticipant(topic.Participants, topic.User, "agent")
				if participant != nil {
					log = log.Record("participant", participant.ID)
					chatTopic := gcloudcx.ConversationChatMessageTopic{}.TopicFor(topic.Conversation)
					log.Infof("User's Participant %s state: %s", participant, participant.State)
					switch participant.State {
					case "alerting": // Now we need to "answer" the participant, i.e. turn them connected
						if channel.IsSubscribed(context, chatTopic) {
							continue
						}
						log.Infof("Subscribing to Conversation %s", topic.Conversation)
						_, err := channel.Subscribe(context, gcloudcx.ConversationChatMessageTopic{}.TopicFor(topic.Conversation))
						if err != nil {
							log.Errorf("Failed to subscribe to topic: %s", topic.Name, err)
							continue
						}

						log.Infof("Setting Participant %s state to %s", participant, "connected")
						err = participant.UpdateState(context, topic.Conversation, "connected")
						if err != nil {
							log.Errorf("Failed to set Participant %s state to: %s", participant, "connected", err)
							continue
						}
					case "disconnected": // Finally, if we need tp wrap up the chat, let's do it
						if !channel.IsSubscribed(context, chatTopic) {
							continue
						}
						if participant.WrapupRequired && participant.Wrapup == nil {
							log.Infof("Wrapping up chat")
							// Once the transfer is initiated, we should "Wrapup" the participant
							//   if needed (queue request a wrapup)
							wrapup := &gcloudcx.Wrapup{Code: "Default Wrap-up Code", Name: "Default Wap-up Code"}
							if err := topic.Conversation.Wrapup(context, participant, wrapup); err != nil {
								log.Errorf("Failed to wrapup Participant %s", participant)
								continue
							}
						}
						if err := channel.Unsubscribe(context, gcloudcx.ConversationChatMessageTopic{}.TopicFor(topic.Conversation)); err != nil {
							log.Errorf("Failed to unscubscribe Participant %s  from topic: %s", participant, gcloudcx.ConversationChatMessageTopic{}.TopicFor(topic.Conversation))
							continue
						}
					}
				}
			case *gcloudcx.ConversationChatMessageTopic:
				log = log.Record("conversation", topic.Conversation.ID)
				log.Infof("Conversation: %s, BodyType: %s, Body: %s, sender: %s", topic.Conversation, topic.BodyType, topic.Body, topic.Sender)
				if topic.Type == "message" && topic.BodyType == "standard" { // remove the noise...
					// We need a full conversation object, so we can operate on it
					conversation, err := gcloudcx.Fetch[gcloudcx.ConversationChat](context, client, topic.Conversation)
					if err != nil {
						log.Errorf("Failed to retrieve a Conversation for ID %s", topic.Conversation, err)
						continue
					}
					participant := findParticipant(conversation.Participants, config.User, "agent")
					if participant == nil {
						log.Debugf("%s is not one of the participants of this conversation", config.User)
						continue
					}
					log = log.Record("participant", participant.ID)

					// skip the agent
					if participant.IsMember("chat", topic.Sender) {
						log.Debugf("%s is the sender of this Notification Topic, nothing to do", config.User)
						continue
					}
					if len(participant.Chats) == 0 {
						log.Warnf("Participant's chat id does not exist yet, skipping")
						continue
					}
					// Pretend the Chat Bot is typing... (whereis it is thinking... isn't it?)
					log.Record("chat", participant.Chats[0]).Debugf("The agent is now typing")
					err = conversation.SetTyping(context, participant.Chats[0])
					if err != nil {
						log.Errorf("Failed to send Typing to Chat Member", err)
					}

					// Send stuff to Matt's Google Dialog Flow webservice
					log.Infof("Participant %s, Sending %s Body to Google: %s", participant, topic.BodyType, topic.Body)
					response := struct {
						Intent          string  `json:"intent"`
						Confidence      float64 `json:"confidence"`
						Fulfillment     string  `json:"fulfillmentmessage"`
						EndConversation bool    `json:"end_conversation"`
					}{EndConversation: false}
					if _, err = request.Send(&request.Options{
						URL: config.BotURL,
						Payload: struct {
							Message string `json:"message"`
						}{
							Message: topic.Body,
						},
						Logger: log,
					},
						&response); err != nil {
						log.Errorf("Failed to send text to Google", err)
						continue
					}
					log.Record("response", response).Debugf("Received: %s", response.Fulfillment)
					if err = conversation.Post(context, participant.Chats[0], response.Fulfillment); err != nil {
						log.Errorf("Failed to send Text to Chat Member", err)
					}
					switch {
					case response.EndConversation:
						log.Infof("Disconnecting Participant %s", participant)
						if err := conversation.Disconnect(context, participant); err != nil {
							log.Errorf("Failed to Wrapup Participant %s", &participant, err)
							continue
						}
					case "agenttransfer" == strings.ToLower(response.Intent):
						log.Infof("Transferring Participant %s to Queue %s", participant, config.AgentQueue)
						log.Record("queue", config.AgentQueue).Debugf("Agent Queue: %s", config.AgentQueue)
						if err := conversation.Transfer(context, participant, config.AgentQueue); err != nil {
							log.Errorf("Failed to Transfer Participant %s to Queue %s", &participant, config.AgentQueue, err)
							continue
						}
					}
				}
			case *gcloudcx.UserPresenceTopic:
				log.Infof("User %s, Presence: %s", topic.User, topic.Presence)
			default:
				log.Warnf("Unknown topic: %s", topic)
			}
		case <-time.After(30 * time.Second): // This timer makes sure the for loop does not execute too quickly when no topic is received for a while
			// log.Debugf("Nothing in the last 30 seconds")
		}
	}
}
