package main

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
)

// AppConfig describes An Application Configuration
type AppConfig struct {
	// WebRootPath contains the root path to prepend to URL when redirecting or in the web pages
	WebRootPath string

	// BotURL is the Chat Bot's url
	BotURL *url.URL

	// BotQueue is the Queue used for sending the initial chat to
	BotQueue *gcloudcx.Queue

	// AgentQueue is the Queue used when the customer wants to talk to an agent
	AgentQueue *gcloudcx.Queue

	// User is the currently Logged in User
	User *gcloudcx.User

	// NotificationChannel is the channel used to receive notifications from GCloud
	NotificationChannel *gcloudcx.NotificationChannel

	// Logger is the Logger for the configuration
	Logger *logger.Logger
}

type key int

// AppConfigKey is the key to store AppConfig in context.Context
const AppConfigKey key = iota

// ToContext store this AppConfig in the given context
func (config *AppConfig) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, AppConfigKey, config)
}

// AppConfigFromContext retrieves an AppConfig from the given context
func AppConfigFromContext(context context.Context) (*AppConfig, error) {
	value := context.Value(AppConfigKey)
	if value == nil {
		return nil, errors.ArgumentMissing.With("AppConfig")
	}
	if config, ok := value.(*AppConfig); ok {
		return config, nil
	}
	return nil, errors.ArgumentInvalid.With("AppConfig", value)
}

// HttpHandler wraps this AppConfig into an http.Handler
func (config *AppConfig) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(config.ToContext(r.Context())))
		})
	}
}

// Initialize configures this AppConfig by calling GCloud as needed
func (config *AppConfig) Initialize(context context.Context, client *gcloudcx.Client) (err error) {
	var correlationID string
	config.Logger = client.Logger.Topic("config")
	log := config.Logger.Scope("initialize")
	if config.AgentQueue == nil {
		return errors.ArgumentMissing.With("Queue")
	}
	match := func(queuename string) func(queue gcloudcx.Queue) bool {
		return func(queue gcloudcx.Queue) bool {
			return strings.EqualFold(queue.Name, queuename)
		}
	}
	if len(config.AgentQueue.ID) == 0 {
		config.AgentQueue, correlationID, err = gcloudcx.FetchBy(context, client, match(config.AgentQueue.Name))
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the Agent Queue %s", config.AgentQueue.Name)
		}
	}
	if config.BotQueue == nil {
		return errors.ArgumentMissing.With("Bot Queue")
	}
	if len(config.BotQueue.ID) == 0 {
		config.BotQueue, correlationID, err = gcloudcx.FetchBy(context, client, match(config.BotQueue.Name))
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the Bot Queue %s", config.BotQueue.Name)
		}
	}

	config.User, err = client.GetMyUser(context)
	if err != nil {
		log.Errorf("Failed to retrieve my User", err)
		return
	}
	log.Infof("Current User: %s", config.User)

	config.NotificationChannel, correlationID, err = client.CreateNotificationChannel(context)
	if err != nil {
		log.Record("genesys-correlation", correlationID).Errorf("Failed to create a notification channel", err)
		return
	}

	topics, correlationID, err := config.NotificationChannel.Subscribe(
		context,
		gcloudcx.UserPresenceTopic{}.With(config.User),
		gcloudcx.UserConversationChatTopic{}.With(config.User),
	)
	if err != nil {
		log.Record("genesys-correlation", correlationID).Errorf("Failed to subscribe to topics", err)
		return
	}
	log.Infof("Subscribed to topics: [%v]", topics)

	// Call the Gcloud Notification Topic loop
	go MessageLoop(config, client)

	return nil
}

// Reset resets the config
func (config *AppConfig) Reset(context context.Context) (err error) {
	log := config.Logger.Scope("reset")

	if config.NotificationChannel != nil {
		log.Debugf("Closing Notification Channel %s", config.NotificationChannel)
		correlationID, err := config.NotificationChannel.Close(context)
		if err != nil {
			config.NotificationChannel.Logger.Record("genesys-correlation", correlationID).Errorf("Failed to close the notification channel %s", config.NotificationChannel, err)
			return err
		}
		log.Infof("Closed Notification Channel %s", config.NotificationChannel)
		config.NotificationChannel = nil
	}
	config.User = nil
	log.Infof("Config is reset")
	return
}
