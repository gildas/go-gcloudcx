package main

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
)

// findParticipant finds a participant after its user id and purpose
func findParticipant(participants []*purecloud.Participant, user *purecloud.User, purpose string) *purecloud.Participant {
	for _, participant := range participants {
		if participant.Purpose == purpose && participant.User != nil && strings.Compare(user.ID, participant.User.ID) == 0 {
			return participant
		}
	}
	return nil
}

// MainHandler is the main webpage. It displays some login info and a WebChat widget
func MainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Topic("route").Scope("main")
		appConfig, err := AppConfigFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the Application Configuration", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Initialize data for the Main Page Template
		viewData := struct {
			Region         string
			DeploymentID   string
			OrganizationID string
			AgentQueue     *purecloud.Queue
			BotQueue       *purecloud.Queue
			BotQueueID     string
			User           *purecloud.User
			ChannelID      string
			WebsocketURL   string
			WebRootPath    string
			LoggedIn       bool
		}{
			WebRootPath: appConfig.WebRootPath,
			LoggedIn:    client.IsAuthorized(),
		}

		// We can use the client only if the agent is logged in...
		if viewData.LoggedIn {
			viewData.Region         = client.Region
			viewData.DeploymentID   = client.DeploymentID
			viewData.OrganizationID = client.Organization.ID
			viewData.AgentQueue     = appConfig.AgentQueue
			viewData.BotQueue       = appConfig.BotQueue

			user, err := client.GetMyUser()
			if err != nil {
				log.Errorf("Failed to retrieve my User", err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
			viewData.User = user

			channel, err := client.CreateNotificationChannel()
			if err != nil {
				log.Errorf("Failed to create a notification channel", err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
			viewData.ChannelID    = channel.ID
			viewData.WebsocketURL = channel.ConnectURL.String()

			topics, err := channel.Subscribe(
				purecloud.UserPresenceTopic{}.TopicFor(user),
				purecloud.UserConversationChatTopic{}.TopicFor(user),
			)
			if err != nil {
				log.Errorf("Failed to subscribe to topics", err)
				core.RespondWithError(w, http.StatusServiceUnavailable, err)
				return
			}
			log.Infof("Subscribed to topics: [%s]", strings.Join(topics, ","))

			go func() {
				log = log.Topic("topic").Scope("process")
				// Processing Received NotificationTopic by reading the chan
				// We do this in a non-blocking way with a timeout to not loop too fast
				// TODO: Add a chan to stop the goroutine that would be written when the route /logout is executed
				for {
					select {
					case receivedTopic := <-channel.TopicReceived:
						log.Debugf("Received topic: %s", receivedTopic)
						switch topic := receivedTopic.(type) {
						case *purecloud.UserConversationChatTopic:
							log = log.Record("user", topic.User.ID).Record("conversation", topic.Conversation.ID)
							log.Infof("User %s, Conversation: %s (state: %s)", topic.User, topic.Conversation, topic.Conversation.State)
							for i, participant := range topic.Participants {
								log.Infof("  Participant #%d: id=%s, name=%s, purpose=%s, state=%s", i, participant.ID, participant.Name, participant.Purpose, participant.State)
							}
							participant := findParticipant(topic.Participants, topic.User, "agent")
							if participant != nil {
								log = log.Record("participant", participant.ID)
								chatTopic := purecloud.ConversationChatMessageTopic{}.TopicFor(topic.Conversation)
								log.Infof("User's Participant %s state: %s", participant, participant.State)
								switch participant.State {
								case "alerting": // Now we need to "answer" the participant, i.e. turn them connected
									if channel.IsSubscribed(chatTopic) {
										continue
									}
									log.Infof("Subscribing to Conversation %s", topic.Conversation)
									_, err := channel.Subscribe(purecloud.ConversationChatMessageTopic{}.TopicFor(topic.Conversation))
									if err != nil {
										log.Errorf("Failed to subscribe to topic: %s", topic.Name, err)
										continue
									}

									log.Infof("Setting Participant %s state to %s", participant, "connected")
									err = topic.Conversation.SetStateParticipant(participant, "connected")
									if err != nil {
										log.Errorf("Failed to set Participant %s state to: %s", participant, "connected", err)
										continue
									}
								case "disconnected": // Finally, if we need tp wrap up the chat, let's do it
									if !channel.IsSubscribed(chatTopic) {
										continue
									}
									if participant.WrapupRequired && participant.Wrapup == nil {
										log.Infof("Wrapping up chat")
										// Once the transfer is initiated, we should "Wrapup" the participant
										//   if needed (queue request a wrapup)
										wrapup := &purecloud.Wrapup{Code: "Default Wrap-up Code", Name: "Default Wap-up Code"}
										if err := topic.Conversation.WrapupParticipant(participant, wrapup); err != nil {
											log.Errorf("Failed to wrapup Participant %s", participant)
											continue
										}
									}
									if err := channel.Unsubscribe(purecloud.ConversationChatMessageTopic{}.TopicFor(topic.Conversation)); err != nil {
										log.Errorf("Failed to unscubscribe Participant %s  from topic: %s", participant, purecloud.ConversationChatMessageTopic{}.TopicFor(topic.Conversation))
										continue
									}
								}
							}
						case *purecloud.ConversationChatMessageTopic:
							log = log.Record("conversation", topic.Conversation.ID)
							log.Infof("Conversation: %s, BodyType: %s, Body: %s, sender: %s", topic.Conversation, topic.BodyType, topic.Body, topic.Sender)
							if topic.Type == "message" && topic.BodyType == "standard" { // remove the noise...
								// We need a real conversation object, so we can operate on it
								err = topic.Conversation.GetMyself()
								if err != nil {
									log.Errorf("Failed to retreive a Conversation for ID %s", topic.Conversation, err)
									continue
								}
								participant := findParticipant(topic.Conversation.Participants, user, "agent")
								if participant != nil {
									log = log.Record("participant", participant.ID)
									switch {
									case strings.Contains(topic.Body, "stop"): // the agent wants to disconnect
										log.Infof("Disconnecting Participant %s", participant)
										if err := topic.Conversation.DisconnectParticipant(participant); err != nil {
											log.Errorf("Failed to Wrapup Participant %s", &participant, err)
											continue
										}
									case strings.Contains(topic.Body, "agent"):
										log.Infof("Transferring Participant %s to Queue %s", participant, appConfig.AgentQueue)
										log.Record("queue", appConfig.AgentQueue).Debugf("Agent Queue: %s", appConfig.AgentQueue)
										if err := topic.Conversation.TransferParticipant(participant, appConfig.AgentQueue); err != nil {
											log.Errorf("Failed to Transfer Participant %s to Queue %s", &participant, appConfig.AgentQueue, err)
											continue
										}
									default: // send the message to the Chat Bot (customer side only)
										if !participant.IsMember("chat", topic.Sender) {
											log.Infof("Participant %s, Sending %s Body to Google: %s", participant, topic.BodyType, topic.Body)

											err = topic.Conversation.SetTyping(participant.Chats[0])
											if err != nil {
												log.Errorf("Failed to send Typing to Chat Member", err)
											}

											// Send stuff to Google
											googleBotURL, _ := url.Parse("https://newpod-gaap.live.genesys.com/MattGDF/")
											response := struct {
												Intent          string            `json:"intent"`
												Confidence      int               `json:"confidence"`
												FulfillmentText string            `json:"fulfillmenttext"`
												Entities        map[string]string `json:"entities"`
											}{}
											_, err := core.SendRequest(&core.RequestOptions{
												URL: googleBotURL,
												Payload: struct {
													Message string `json:"message"`
												}{
													Message: topic.Body,
												},
											},
												&response)
											if err != nil {
												log.Errorf("Failed to send text to Google", err)
												continue
											}

											log.Record("response", response).Debugf("Received: %s", response.FulfillmentText)
											err = topic.Conversation.Post(participant.Chats[0], response.FulfillmentText)
											if err != nil {
												log.Errorf("Failed to send Text to Chat Member", err)
											}
										}
									}
								} else {
									log.Warnf("Failed to find Agent Participant in Conversation")
								}
							}
						case *purecloud.UserPresenceTopic:
							log.Infof("User %s, Presence: %s", topic.User, topic.Presence)
						default:
							log.Warnf("Unknown topic: %s", topic)
						}
					case <-time.After(30 * time.Second): // This timer makes sure the for loop does not execute too quickly when no topic is received for a while
						// log.Debugf("Nothing in the last 30 seconds")
					}
				}
			}()
		}

		log.Infof(`Rendering page "page_main"`)
		pageTemplate, err := template.ParseFiles("page_main.html")
		if err != nil {
			log.Errorf(`Failed to parse page "page_main"`, err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		err = pageTemplate.Execute(w, viewData)
		if err != nil {
			log.Errorf(`Failed to render page "page_main"`, err)
		}
	})
}
