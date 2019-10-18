package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
	"github.com/gorilla/mux"
)

// Log is the application Logger
var Log *logger.Logger

// Client is the PureCloud Client
var Client *purecloud.Client

// The Queue to transfer to
var Queue *purecloud.Queue

func loggedInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Scope("logged_in")
		log.Infof("Redirecting to /")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}

func mainRouteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Scope("main")

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		user, err := client.GetMyUser()
		if err != nil {
			log.Errorf("Failed to retrieve my User", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		channel, err := client.CreateNotificationChannel()
		if err != nil {
			log.Errorf("Failed to create a notification channel", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

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
			// TODO: Add a chan to stop the goroutine
			for {
				select {
				case receivedTopic := <-channel.TopicReceived:
					log.Infof("Received topic: %s", receivedTopic)
					switch topic := receivedTopic.(type) {
					case *purecloud.UserConversationChatTopic:
						log.Infof("User %s, Conversation: %s", topic.UserID, topic.ConversationID)
						// TODO: Matt => What is that connected variable in index.html?!?
						if len(topic.Participants) >= 4 && len(topic.Participants[3].ID) != 0 {
							log.Infof("Subscribing to Conversation %s", topic.ConversationID)
							conversation := purecloud.Conversation{ID: topic.ConversationID}
							_, err := channel.Subscribe(purecloud.ConversationChatMessageTopic{}.TopicFor(conversation))
							if err != nil {
								log.Errorf("Failed to subscribe to topic: %s", topic.Name, err)
								continue
							}
							// Now we need to "answer" the participant, i.e. turn them connected
							participant := conversation.Participants[3]
							err = conversation.SetStateParticipant(participant, "connected")
							if err != nil {
								log.Errorf("Failed to set Participant %s state to: %s", participant, "connected", err)
								continue
							}
						}
					case *purecloud.ConversationChatMessageTopic:
						log.Infof("Conversation: %s, BodyType: %s, Body: %s", topic.ConversationID, topic.BodyType, topic.Body)
						// We need a real conversation object, so we can operate on it
						conversation, err := topic.GetClient().GetConversation(topic.ConversationID)
						if err != nil {
							log.Errorf("Failed to retreive a Conversation for ID %s", topic.ConversationID, err)
							continue
						}
						participant := conversation.Participants[3]
						if strings.Contains(topic.Body, "stop") { // disconnect
							if err := conversation.DisconnectParticipant(participant); err != nil {
								log.Errorf("Failed to Wrapup Participant %s", &participant, err)
								continue
							}
						} else if strings.Contains(topic.Body, "agent") { // transfer
							if err := conversation.TransferParticipant(participant, Queue); err != nil {
								log.Errorf("Failed to Transfer Participant %s to Queue %s", &participant, Queue, err)
								continue
							}
							// Once the transfer is initiated, we should "Wrapup" the participant
							//   if needed (queue request a wrapup)
							wrapup := &purecloud.Wrapup{Code: "Default Wrap-up Code", Name: "Default Wap-up Code"}
							if err := conversation.WrapupParticipant(participant, wrapup); err != nil {
								log.Errorf("Failed to wrapup Partitipant %s", participant)
								continue
							}
						} else {
							log.Infof("Sending %s Body to Google: %s", topic.BodyType, topic.Body)
							// Send stuff to Google
						}
					case *purecloud.UserPresenceTopic:
						log.Infof("User %s, Presence: %s", topic.UserID, topic.Presence)
					default:
						log.Warnf("Unknown topic: %s", topic)
					}
				case <-time.After(30 * time.Second):
					log.Debugf("Nothing in the last 30 seconds")
				}
			}
		}()

		core.RespondWithJSON(w, http.StatusOK, struct {
			UserName     string `json:"user"`
			ChannelID    string `json:"channelId"`
			WebsocketURL string `json:"websocketUrl"`
		}{
			UserName:     user.Name,
			ChannelID:    channel.ID,
			WebsocketURL: channel.ConnectURL.String(),
		})
	})
}

func main() {
	var (
		region       = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the PureCloud Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the PureCloud Client ID for authentication")
		secret       = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the PureCloud Client Secret for authentication")
		deploymentID = flag.String("deploymentid", core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""), "the PureCloud Application Deployment ID")
		redirectRoot = flag.String("redirecturi", core.GetEnvAsString("PURECLOUD_REDIRECTURI", ""), "The root uri to give to PureCloud as a Redirect URI")
		queueID      = flag.String("queue", core.GetEnvAsString("PURECLOUD_QUEUE", ""), "The queue to transfer to")
		port         = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
	)
	flag.Parse()

	Log = logger.Create("ChatBot_Example")

	Log.Infof("redirect root: %s", *redirectRoot)
	if len(*redirectRoot) == 0 {
		*redirectRoot = fmt.Sprintf("http://localhost:%d", *port)
	}
	redirectURL, err := url.Parse(*redirectRoot + "/token")
	if err != nil {
		Log.Fatalf("Invalid Redirect URL: %s/token", *redirectRoot, err)
		os.Exit(-1)
	}
	Log.Infof("Make sure your PureCloud OAUTH accepts redirects to: %s", redirectURL.String())

	Client = purecloud.New(purecloud.ClientOptions{
		Region:       *region,
		DeploymentID: *deploymentID,
		Logger:       Log,
	}).SetAuthorizationGrant(&purecloud.AuthorizationCodeGrant{
		ClientID:    *clientID,
		Secret:      *secret,
		RedirectURL: redirectURL,
	})

	// TODO: Make this better... Too Simple for now
	Queue = &purecloud.Queue{ID: *queueID}

	// Create the HTTP Incoming Request Router
	router := mux.NewRouter().StrictSlash(true)
	// This route actually performs login the user using the grant of the purecloud.Client
	//   Upon success, your route httpHandler is called
	router.Methods("GET").Path("/token").Handler(Log.HttpHandler()(Client.LoginHandler()(loggedInHandler())))

	// This route performs your actions, but makes sure the client is authorized,
	//   if authorized, your route http.Handler is called
	//   otherwise, the AuthorizeHandler will redirect the user to the PureCloud Login page
	//   that will end up with the grant.RedirectURL defined ealier
	router.Methods("GET").Path("/").Handler(Log.HttpHandler()(Client.AuthorizeHandler()(mainRouteHandler())))

	WebServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", *port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
		//ErrorLog:     Log,
	}

	// Starting the server
	go func() {
		log := Log.Topic("webserver").Scope("run")

		log.Infof("Starting WEB server on port %d", *port)
		if err := WebServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Fatalf("Failed to start the WEB server on port: %d", *port, err)
			}
		}
	}()

	// Accepting shutdowns from SIGINT (^C) and SIGTERM (docker, heroku)
	interruptChannel := make(chan os.Signal, 1)
	exitChannel := make(chan struct{})

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// The go routine that wait for cleaning stuff when exiting
	go func() {
		sig := <-interruptChannel // Block until we have to stop

		context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		Log.Infof("Application is stopping (%+v)", sig)

		// Stopping the WEB server
		Log.Debugf("WEB server is shutting down")
		WebServer.SetKeepAlivesEnabled(false)
		err := WebServer.Shutdown(context)
		if err != nil {
			Log.Errorf("Failed to stop the WEB server", err)
		} else {
			Log.Infof("WEB server is stopped")
		}
		close(exitChannel)
	}()

	<-exitChannel
	os.Exit(0)
}
